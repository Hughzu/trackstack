package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/Hughzu/trackstack/apps/server/internal/contexts/heat/application/ports"
	"github.com/Hughzu/trackstack/apps/server/internal/platform/authcontext"
	"github.com/Hughzu/trackstack/apps/server/internal/platform/timeutil"
)

type createRefillPayload struct {
	Date        string   `json:"date"`
	WeightKg    float64  `json:"weightKg"`
	Bags        int      `json:"bags"`
	Temperature *float64 `json:"temperature"`
}

type errorResponse struct {
	Error string `json:"error"`
}

type RefillHandler struct {
	useCase ports.RefillUseCase
}

func NewRefillHandler(useCase ports.RefillUseCase) *RefillHandler {
	return &RefillHandler{useCase: useCase}
}

func (h *RefillHandler) GetRefills(w http.ResponseWriter, r *http.Request) {
	userID, ok := h.userID(w, r)
	if !ok {
		return
	}

	fromStr := r.URL.Query().Get("from")
	toStr := r.URL.Query().Get("to")

	from, err := timeutil.ParseDate(fromStr)
	if err != nil {
		h.writeJSON(w, http.StatusBadRequest, errorResponse{Error: "Invalid 'from' date"})
		return
	}

	to, err := timeutil.ParseDate(toStr)
	if err != nil {
		h.writeJSON(w, http.StatusBadRequest, errorResponse{Error: "Invalid 'to' date"})
		return
	}

	refills, err := h.useCase.GetRefills(r.Context(), ports.GetRefillsQuery{
		UserID: userID,
		From:   from,
		To:     to,
	})
	if err != nil {
		h.writeError(w, err)
		return
	}

	h.writeJSON(w, http.StatusOK, refills)
}

func (h *RefillHandler) CreateRefill(w http.ResponseWriter, r *http.Request) {
	userID, ok := h.userID(w, r)
	if !ok {
		return
	}

	var payload createRefillPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		h.writeJSON(w, http.StatusBadRequest, errorResponse{Error: "Invalid JSON body"})
		return
	}

	date, err := timeutil.ParseDate(payload.Date)
	if err != nil {
		h.writeJSON(w, http.StatusBadRequest, errorResponse{Error: "Invalid 'date'"})
		return
	}

	refill, err := h.useCase.CreateRefill(r.Context(), ports.CreateRefillCommand{
		UserID:      userID,
		Date:        date,
		WeightKg:    payload.WeightKg,
		Bags:        payload.Bags,
		Temperature: payload.Temperature,
	})
	if err != nil {
		h.writeError(w, err)
		return
	}

	h.writeJSON(w, http.StatusCreated, refill)
}

func (h *RefillHandler) deleteRefill(w http.ResponseWriter, r *http.Request, id string) {
	userID, ok := h.userID(w, r)
	if !ok {
		return
	}

	if id == "" {
		h.writeJSON(w, http.StatusBadRequest, errorResponse{Error: "Missing refill id"})
		return
	}

	deleted, err := h.useCase.DeleteRefill(r.Context(), ports.DeleteRefillCommand{
		UserID: userID,
		ID:     id,
	})
	if err != nil {
		h.writeError(w, err)
		return
	}
	if !deleted {
		h.writeJSON(w, http.StatusNotFound, errorResponse{Error: "Refill not found"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *RefillHandler) writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func (h *RefillHandler) writeError(w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	if errors.Is(err, ports.ErrInvalidInput) {
		status = http.StatusBadRequest
	}

	h.writeJSON(w, status, errorResponse{Error: err.Error()})
}

func (h *RefillHandler) GetDashboard(w http.ResponseWriter, r *http.Request) {
	userID, ok := h.userID(w, r)
	if !ok {
		return
	}

	page := 1
	if pageValue := r.URL.Query().Get("page"); pageValue != "" {
		if parsed, err := strconv.Atoi(pageValue); err == nil && parsed > 0 {
			page = parsed
		}
	}

	limit := 20
	if limitValue := r.URL.Query().Get("limit"); limitValue != "" {
		if parsed, err := strconv.Atoi(limitValue); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	dashboard, err := h.useCase.GetDashboard(r.Context(), ports.GetDashboardQuery{
		UserID: userID,
		Page:   page,
		Limit:  limit,
	})
	if err != nil {
		h.writeError(w, err)
		return
	}

	h.writeJSON(w, http.StatusOK, dashboard)
}

func (h *RefillHandler) userID(w http.ResponseWriter, r *http.Request) (string, bool) {
	userID, ok := authcontext.GetUserID(r.Context())
	if !ok || userID == "" {
		h.writeJSON(w, http.StatusUnauthorized, errorResponse{Error: "Unauthorized"})
		return "", false
	}

	return userID, true
}
