package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (h *CaloriesHandler) RegisterRoutes(r chi.Router) {
	r.Get("/dashboard", h.GetDashboard)
	r.Get("/target", h.GetTarget)
	r.Post("/target", h.UpdateTarget)
	r.Post("/log", h.AddLog)
	r.Delete("/logs/{id}", func(w http.ResponseWriter, r *http.Request) {
		h.DeleteLog(w, r, chi.URLParam(r, "id"))
	})
}
