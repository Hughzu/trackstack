package httptransport

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Hughzu/trackstack/apps/server/internal/modules/heat"
)

type HeatHandler struct {
	svc    *heat.Service
	userID string
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
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")

	resp, err := h.svc.ListRefills(r.Context(), heat.ListRefillsRequest{
		UserID: h.userID,
		From:   from,
		To:     to,
	})
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *HeatHandler) CreateRefill(w http.ResponseWriter, r *http.Request) {
	var payload createRefillPayload
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&payload); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "Invalid JSON body"})
		return
	}

	resp, err := h.svc.CreateRefill(r.Context(), heat.CreateRefillRequest{
		UserID:      h.userID,
		Date:        payload.Date,
		WeightKg:    payload.WeightKg,
		Bags:        payload.Bags,
		Temperature: payload.Temperature,
	})
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, resp)
}

func writeError(w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	if errors.Is(err, heat.ErrInvalidInput) {
		status = http.StatusBadRequest
	}

	writeJSON(w, status, errorResponse{Error: err.Error()})
}
