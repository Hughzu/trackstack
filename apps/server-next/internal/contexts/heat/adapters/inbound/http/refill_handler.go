package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Hughzu/trackstack/apps/server-next/internal/contexts/heat/application/ports"
	"github.com/Hughzu/trackstack/apps/server-next/internal/platform/timeutil"
)

const mockUserID = "8a36e9e2-4b42-4ea2-a397-0a2b441accca"

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
		UserID: mockUserID,
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
		UserID:      mockUserID,
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
	if id == "" {
		h.writeJSON(w, http.StatusBadRequest, errorResponse{Error: "Missing refill id"})
		return
	}

	deleted, err := h.useCase.DeleteRefill(r.Context(), ports.DeleteRefillCommand{
		UserID: mockUserID,
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
