package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (h *RefillHandler) RegisterRoutes(r chi.Router) {
	r.Get("/dashboard", h.GetDashboard)
	r.Get("/refills", h.GetRefills)
	r.Post("/refills", h.CreateRefill)
	r.Delete("/refills/{id}", func(w http.ResponseWriter, r *http.Request) {
		h.deleteRefill(w, r, chi.URLParam(r, "id"))
	})
}
