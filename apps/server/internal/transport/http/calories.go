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

func (h *CaloriesHandler) UpdateTarget(w http.ResponseWriter, r *http.Request) {
	userID, ok := requireAuthUserID(w, r)
	if !ok {
		return
	}
	isJSON := isJSONRequest(r)
	data, ok := readCaloriesPayload(r, isJSON)
	if !ok {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "Invalid JSON body"})
		return
	}

	if data.TargetKcal == nil || data.TargetProtein == nil {
		if !isJSON {
			redirectWithError(w, r, "/calories/settings")
			return
		}
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

	if !isJSON {
		redirect(w, "/calories")
		return
	}

	writeJSON(w, http.StatusOK, target)
}

func (h *CaloriesHandler) AddLog(w http.ResponseWriter, r *http.Request) {
	userID, ok := requireAuthUserID(w, r)
	if !ok {
		return
	}
	isJSON := isJSONRequest(r)
	data, ok := readCaloriesPayload(r, isJSON)
	if !ok {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "Invalid JSON body"})
		return
	}

	if data.Calories == nil || data.Protein == nil {
		if !isJSON {
			redirectWithError(w, r, "/calories")
			return
		}
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

	if !isJSON {
		redirect(w, "/calories")
		return
	}

	writeJSON(w, http.StatusCreated, log)
}

func (h *CaloriesHandler) DeleteLog(w http.ResponseWriter, r *http.Request) {
	userID, ok := requireAuthUserID(w, r)
	if !ok {
		return
	}
	var id string
	if payload, ok := readJSONMap(r); ok {
		if value, ok := payload["id"]; ok {
			id = parseOptionalString(value)
		}
	}

	if id == "" {
		id = r.URL.Query().Get("id")
	}

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

func readCaloriesPayload(r *http.Request, isJSON bool) (caloriesPayload, bool) {
	if isJSON {
		payload, ok := readJSONMap(r)
		if !ok {
			return caloriesPayload{}, false
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
		}, true
	}

	if err := r.ParseForm(); err != nil {
		return caloriesPayload{}, true
	}

	return caloriesPayload{
		Calories:      parseOptionalIntString(r.FormValue("calories")),
		Protein:       parseOptionalIntString(r.FormValue("protein")),
		Carbs:         parseOptionalIntString(r.FormValue("carbs")),
		Fat:           parseOptionalIntString(r.FormValue("fat")),
		Title:         parseOptionalStringPtr(r.FormValue("title")),
		Date:          parseOptionalStringPtr(r.FormValue("date")),
		TargetKcal:    parseOptionalIntString(r.FormValue("targetKcal")),
		TargetProtein: parseOptionalIntString(r.FormValue("targetProtein")),
		TargetCarbs:   parseOptionalIntString(r.FormValue("targetCarbs")),
		TargetFat:     parseOptionalIntString(r.FormValue("targetFat")),
	}, true
}

func readJSONMap(r *http.Request) (map[string]any, bool) {
	var payload map[string]any
	if err := decodeJSON(r, &payload); err != nil {
		return nil, false
	}
	return payload, true
}

func parseOptionalInt(value any) *int {
	switch v := value.(type) {
	case float64:
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

func parseOptionalString(value any) string {
	if text, ok := value.(string); ok {
		return strings.TrimSpace(text)
	}
	return ""
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

func isJSONRequest(r *http.Request) bool {
	return strings.Contains(r.Header.Get("Content-Type"), "application/json")
}

func redirectWithError(w http.ResponseWriter, r *http.Request, fallback string) {
	referrer := strings.TrimSpace(r.Referer())
	if referrer != "" {
		fallback = referrer
	}

	if strings.Contains(fallback, "?") {
		fallback = fallback + "&error=1"
	} else {
		fallback = fallback + "?error=1"
	}

	redirect(w, fallback)
}

func redirect(w http.ResponseWriter, location string) {
	w.Header().Set("Location", location)
	w.WriteHeader(http.StatusSeeOther)
}
