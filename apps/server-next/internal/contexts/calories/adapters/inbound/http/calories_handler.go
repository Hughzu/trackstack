package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/Hughzu/trackstack/apps/server-next/internal/contexts/calories/application/ports"
	"github.com/Hughzu/trackstack/apps/server-next/internal/platform/authcontext"
)

type errorResponse struct {
	Error string `json:"error"`
}

type caloriesPayload struct {
	Calories           *int
	ProteinGrams       *int
	CarbGrams          *int
	FatGrams           *int
	Title              *string
	Date               *string
	Time               *string
	TargetCalories     *int
	TargetProteinGrams *int
	TargetCarbGrams    *int
	TargetFatGrams     *int
}

type CaloriesHandler struct {
	targetUseCase    ports.TargetUseCase
	logUseCase       ports.LogUseCase
	dashboardUseCase ports.DashboardUseCase
}

func NewCaloriesHandler(targetUseCase ports.TargetUseCase, logUseCase ports.LogUseCase, dashboardUseCase ports.DashboardUseCase) *CaloriesHandler {
	return &CaloriesHandler{
		targetUseCase:    targetUseCase,
		logUseCase:       logUseCase,
		dashboardUseCase: dashboardUseCase,
	}
}

func (h *CaloriesHandler) GetTarget(w http.ResponseWriter, r *http.Request) {
	userID, ok := h.userID(w, r)
	if !ok {
		return
	}

	target, err := h.targetUseCase.GetTarget(r.Context(), ports.GetTargetQuery{UserID: userID})
	if err != nil {
		h.writeError(w, err)
		return
	}

	h.writeJSON(w, http.StatusOK, target)
}

func (h *CaloriesHandler) UpdateTarget(w http.ResponseWriter, r *http.Request) {
	userID, ok := h.userID(w, r)
	if !ok {
		return
	}

	data, err := readCaloriesPayload(r)
	if err != nil {
		h.writeJSON(w, http.StatusBadRequest, errorResponse{Error: "Invalid JSON body"})
		return
	}
	if data.TargetCalories == nil || data.TargetProteinGrams == nil {
		h.writeJSON(w, http.StatusBadRequest, errorResponse{Error: "Missing required fields"})
		return
	}

	target, err := h.targetUseCase.UpdateTarget(r.Context(), ports.UpdateTargetCommand{
		UserID:             userID,
		TargetCalories:     *data.TargetCalories,
		TargetProteinGrams: *data.TargetProteinGrams,
		TargetCarbGrams:    data.TargetCarbGrams,
		TargetFatGrams:     data.TargetFatGrams,
	})
	if err != nil {
		h.writeError(w, err)
		return
	}

	h.writeJSON(w, http.StatusOK, target)
}

func (h *CaloriesHandler) AddLog(w http.ResponseWriter, r *http.Request) {
	userID, ok := h.userID(w, r)
	if !ok {
		return
	}

	data, err := readCaloriesPayload(r)
	if err != nil {
		h.writeJSON(w, http.StatusBadRequest, errorResponse{Error: "Invalid JSON body"})
		return
	}
	if data.Calories == nil || data.ProteinGrams == nil {
		h.writeJSON(w, http.StatusBadRequest, errorResponse{Error: "Missing required fields"})
		return
	}

	log, err := h.logUseCase.AddLog(r.Context(), ports.AddLogCommand{
		UserID:       userID,
		Calories:     *data.Calories,
		ProteinGrams: *data.ProteinGrams,
		CarbGrams:    data.CarbGrams,
		FatGrams:     data.FatGrams,
		Title:        data.Title,
		Date:         data.Date,
		Time:         data.Time,
	})
	if err != nil {
		h.writeError(w, err)
		return
	}

	h.writeJSON(w, http.StatusCreated, log)
}

func (h *CaloriesHandler) DeleteLog(w http.ResponseWriter, r *http.Request, id string) {
	userID, ok := h.userID(w, r)
	if !ok {
		return
	}

	if strings.TrimSpace(id) == "" {
		h.writeJSON(w, http.StatusBadRequest, errorResponse{Error: "Missing log id"})
		return
	}

	deleted, err := h.logUseCase.DeleteLog(r.Context(), ports.DeleteLogCommand{
		UserID: userID,
		ID:     id,
	})
	if err != nil {
		h.writeError(w, err)
		return
	}
	if !deleted {
		h.writeJSON(w, http.StatusNotFound, errorResponse{Error: "Log not found"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *CaloriesHandler) GetDashboard(w http.ResponseWriter, r *http.Request) {
	userID, ok := h.userID(w, r)
	if !ok {
		return
	}

	recentLimit := 8
	if limitValue := r.URL.Query().Get("recentLimit"); limitValue != "" {
		if parsed, err := strconv.Atoi(limitValue); err == nil && parsed > 0 {
			recentLimit = parsed
		}
	}

	logsLimit := 50
	if limitValue := r.URL.Query().Get("logsLimit"); limitValue != "" {
		if parsed, err := strconv.Atoi(limitValue); err == nil && parsed > 0 {
			logsLimit = parsed
		}
	}

	dashboard, err := h.dashboardUseCase.GetDashboard(r.Context(), ports.GetDashboardQuery{
		UserID:      userID,
		RecentLimit: recentLimit,
		LogsLimit:   logsLimit,
	})
	if err != nil {
		h.writeError(w, err)
		return
	}

	h.writeJSON(w, http.StatusOK, dashboard)
}

func (h *CaloriesHandler) userID(w http.ResponseWriter, r *http.Request) (string, bool) {
	userID, ok := authcontext.GetUserID(r.Context())
	if !ok || userID == "" {
		h.writeJSON(w, http.StatusUnauthorized, errorResponse{Error: "Unauthorized"})
		return "", false
	}

	return userID, true
}

func readCaloriesPayload(r *http.Request) (caloriesPayload, error) {
	var payload map[string]any
	if err := decodeJSON(r, &payload); err != nil {
		return caloriesPayload{}, err
	}

	return caloriesPayload{
		Calories:           parseOptionalInt(payload["calories"]),
		ProteinGrams:       parseOptionalInt(payload["proteinGrams"]),
		CarbGrams:          parseOptionalInt(payload["carbGrams"]),
		FatGrams:           parseOptionalInt(payload["fatGrams"]),
		Title:              parseOptionalStringPtr(payload["title"]),
		Date:               parseOptionalStringPtr(payload["date"]),
		Time:               parseOptionalStringPtr(payload["time"]),
		TargetCalories:     parseOptionalInt(payload["targetCalories"]),
		TargetProteinGrams: parseOptionalInt(payload["targetProteinGrams"]),
		TargetCarbGrams:    parseOptionalInt(payload["targetCarbGrams"]),
		TargetFatGrams:     parseOptionalInt(payload["targetFatGrams"]),
	}, nil
}

func decodeJSON(r *http.Request, dst any) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(dst)
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

func (h *CaloriesHandler) writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func (h *CaloriesHandler) writeError(w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	if errors.Is(err, ports.ErrInvalidInput) {
		status = http.StatusBadRequest
	}

	h.writeJSON(w, status, errorResponse{Error: err.Error()})
}
