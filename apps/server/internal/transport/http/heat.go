package httptransport

import (
	"errors"
	"net/http"

	"github.com/Hughzu/trackstack/apps/server/internal/modules/heat"
)

type HeatHandler struct {
	svc *heat.Service
}

type createRefillPayload struct {
	Date        string   `json:"date"`
	WeightKg    float64  `json:"weightKg"`
	Bags        int      `json:"bags"`
	Temperature *float64 `json:"temperature"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func (h *HeatHandler) ListRefills(w http.ResponseWriter, r *http.Request) {
	userID, ok := requireAuthUserID(w, r)
	if !ok {
		return
	}
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")

	refills, err := h.svc.ListRefills(r.Context(), heat.ListRefillsRequest{
		UserID: userID,
		From:   from,
		To:     to,
	})
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, refills)
}

func (h *HeatHandler) GetDashboard(w http.ResponseWriter, r *http.Request) {
	userID, ok := requireAuthUserID(w, r)
	if !ok {
		return
	}

	dashboard, err := h.svc.GetDashboard(r.Context(), heat.GetDashboardRequest{
		UserID: userID,
	})
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, dashboard)
}

func (h *HeatHandler) CreateRefill(w http.ResponseWriter, r *http.Request) {
	userID, ok := requireAuthUserID(w, r)
	if !ok {
		return
	}
	var payload createRefillPayload
	if err := decodeJSON(r, &payload); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "Invalid JSON body"})
		return
	}

	refill, err := h.svc.CreateRefill(r.Context(), heat.CreateRefillRequest{
		UserID:      userID,
		Date:        payload.Date,
		WeightKg:    payload.WeightKg,
		Bags:        payload.Bags,
		Temperature: payload.Temperature,
	})
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, refill)
}

func (h *HeatHandler) DeleteRefill(w http.ResponseWriter, r *http.Request) {
	userID, ok := requireAuthUserID(w, r)
	if !ok {
		return
	}
	id := r.URL.Query().Get("id")

	if id == "" {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "Missing refill id"})
		return
	}

	deleted, err := h.svc.DeleteRefill(r.Context(), heat.DeleteRefillRequest{
		UserID: userID,
		ID:     id,
	})
	if err != nil {
		writeError(w, err)
		return
	}
	if !deleted {
		writeJSON(w, http.StatusNotFound, errorResponse{Error: "Refill not found"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func writeError(w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	if errors.Is(err, heat.ErrInvalidInput) {
		status = http.StatusBadRequest
	}

	writeJSON(w, status, errorResponse{Error: err.Error()})
}
