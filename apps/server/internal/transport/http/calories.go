package httptransport

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/Hughzu/trackstack/apps/server/internal/modules/calories"
)

type CaloriesHandler struct {
	svc *calories.Service
}

func (h *CaloriesHandler) GetTarget(w http.ResponseWriter, r *http.Request) {
	userID, ok := requireAuthUserID(w, r)
	if !ok {
		return
	}
	target, err := h.svc.GetTarget(r.Context(), calories.GetTargetRequest{UserID: userID})
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "Server Error"})
		return
	}

	writeJSON(w, http.StatusOK, target)
}

func (h *CaloriesHandler) GetDashboard(w http.ResponseWriter, r *http.Request) {
	userID, ok := requireAuthUserID(w, r)
	if !ok {
		return
	}

	recentLimit := 8
	if limitStr := r.URL.Query().Get("recentLimit"); limitStr != "" {
		if val, err := strconv.Atoi(limitStr); err == nil && val > 0 {
			recentLimit = val
		}
	}

	logsLimit := 50
	if limitStr := r.URL.Query().Get("logsLimit"); limitStr != "" {
		if val, err := strconv.Atoi(limitStr); err == nil && val > 0 {
			logsLimit = val
		}
	}

	dashboard, err := h.svc.GetDashboard(r.Context(), calories.GetDashboardRequest{
		UserID:      userID,
		RecentLimit: recentLimit,
		LogsLimit:   logsLimit,
	})
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "Server Error"})
		return
	}

	writeJSON(w, http.StatusOK, dashboard)
}

func (h *CaloriesHandler) UpdateTarget(w http.ResponseWriter, r *http.Request) {
	userID, ok := requireAuthUserID(w, r)
	if !ok {
		return
	}
	data, err := readCaloriesPayload(r)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "Invalid JSON body"})
		return
	}

	if data.TargetKcal == nil || data.TargetProtein == nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "Missing required fields"})
		return
	}

	target, err := h.svc.UpdateTarget(r.Context(), calories.UpdateTargetRequest{
		UserID:         userID,
		TargetKcal:     data.TargetKcal,
		TargetProteinG: data.TargetProtein,
		TargetCarbsG:   data.TargetCarbs,
		TargetFatG:     data.TargetFat,
	})
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "Server Error"})
		return
	}

	writeJSON(w, http.StatusOK, target)
}

func (h *CaloriesHandler) AddLog(w http.ResponseWriter, r *http.Request) {
	userID, ok := requireAuthUserID(w, r)
	if !ok {
		return
	}
	data, err := readCaloriesPayload(r)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "Invalid JSON body"})
		return
	}

	if data.Calories == nil || data.Protein == nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "Missing required fields"})
		return
	}

	log, err := h.svc.AddLog(r.Context(), calories.AddLogRequest{
		UserID:   userID,
		Calories: data.Calories,
		ProteinG: data.Protein,
		CarbsG:   data.Carbs,
		FatG:     data.Fat,
		Title:    data.Title,
		Date:     data.Date,
		Time:     nil,
	})
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "Server Error"})
		return
	}

	writeJSON(w, http.StatusCreated, log)
}

func (h *CaloriesHandler) DeleteLog(w http.ResponseWriter, r *http.Request) {
	userID, ok := requireAuthUserID(w, r)
	if !ok {
		return
	}
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "Missing log id"})
		return
	}

	deleted, err := h.svc.DeleteLog(r.Context(), calories.DeleteLogRequest{
		UserID: userID,
		ID:     id,
	})
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "Server Error"})
		return
	}
	if !deleted {
		writeJSON(w, http.StatusNotFound, errorResponse{Error: "Log not found"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

type caloriesPayload struct {
	Calories      *int
	Protein       *int
	Carbs         *int
	Fat           *int
	Title         *string
	Date          *string
	TargetKcal    *int
	TargetProtein *int
	TargetCarbs   *int
	TargetFat     *int
}

func readCaloriesPayload(r *http.Request) (caloriesPayload, error) {
	var payload map[string]any
	if err := decodeJSON(r, &payload); err != nil {
		return caloriesPayload{}, err
	}

	return caloriesPayload{
		Calories:      parseOptionalInt(payload["calories"]),
		Protein:       parseOptionalInt(payload["protein"]),
		Carbs:         parseOptionalInt(payload["carbs"]),
		Fat:           parseOptionalInt(payload["fat"]),
		Title:         parseOptionalStringPtr(payload["title"]),
		Date:          parseOptionalStringPtr(payload["date"]),
		TargetKcal:    parseOptionalInt(payload["targetKcal"]),
		TargetProtein: parseOptionalInt(payload["targetProtein"]),
		TargetCarbs:   parseOptionalInt(payload["targetCarbs"]),
		TargetFat:     parseOptionalInt(payload["targetFat"]),
	}, nil
}

func parseOptionalInt(value any) *int {
	switch v := value.(type) {
	case float64:
		if v != float64(int(v)) {
			return nil
		}
		parsed := int(v)
		return &parsed
	case string:
		return parseOptionalIntString(v)
	default:
		return nil
	}
}

func parseOptionalIntString(value string) *int {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil
	}
	parsed, err := strconv.Atoi(trimmed)
	if err != nil {
		return nil
	}
	return &parsed
}

func parseOptionalStringPtr(value any) *string {
	if text, ok := value.(string); ok {
		trimmed := strings.TrimSpace(text)
		if trimmed == "" {
			return nil
		}
		return &trimmed
	}
	return nil
}
