package http

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Hughzu/trackstack/apps/server-next/internal/contexts/heat/application/ports"
)

type RefillResponse struct {
	ID        string `json:"id"`
	Amount    int    `json:"amount"`
	CreatedAt string `json:"created_at"`
}

type RefillHandler struct {
	useCase ports.RefillUseCase
}

func NewRefillHandler(useCase ports.RefillUseCase) *RefillHandler {
	return &RefillHandler{useCase: useCase}
}

func (h *RefillHandler) GetRefills(w http.ResponseWriter, r *http.Request) {

	mockUserId := "8a36e9e2-4b42-4ea2-a397-0a2b441accca"

	fromStr := r.URL.Query().Get("from")
	toStr := r.URL.Query().Get("to")

	from, err := time.Parse(time.RFC3339, fromStr)
	if err != nil {
		http.Error(w, "invalid 'from' format (expected RFC3339)", http.StatusBadRequest)
		return
	}

	to, err := time.Parse(time.RFC3339, toStr)
	if err != nil {
		http.Error(w, "invalid 'to' format (expected RFC3339)", http.StatusBadRequest)
		return
	}

	refills, err := h.useCase.GetRefills(r.Context(), mockUserId, from, to)
	if err != nil {
		http.Error(w, "failed to get refills", http.StatusInternalServerError)
		return
	}

	var response []RefillResponse

	for _, refill := range refills {
		response = append(response, RefillResponse{
			ID:        refill.ID,
			Amount:    refill.Amount,
			CreatedAt: refill.CreatedAt.Format(time.RFC3339),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
