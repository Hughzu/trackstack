package http

import "github.com/go-chi/chi/v5"

func (h *RefillHandler) RegisterRoutes(r chi.Router) {
	r.Get("/refills", h.GetRefills)
}
