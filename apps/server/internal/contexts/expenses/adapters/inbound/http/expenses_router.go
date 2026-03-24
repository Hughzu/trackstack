package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (h *ExpensesHandler) RegisterRoutes(r chi.Router) {
	r.Get("/settings", h.GetSettings)
	r.Post("/settings", h.UpdateSettings)
	r.Get("/sheet/current", h.GetDashboard)
	r.Post("/entries", h.AddEntry)
	r.Delete("/entries/{id}", func(w http.ResponseWriter, r *http.Request) {
		h.DeleteEntry(w, r, chi.URLParam(r, "id"))
	})
	r.Post("/checklists", h.UpsertChecklist)
	r.Delete("/checklists/{id}", func(w http.ResponseWriter, r *http.Request) {
		h.DeleteChecklist(w, r, chi.URLParam(r, "id"))
	})
	r.Post("/checklists/complete", h.CompleteChecklistItem)
	r.Post("/recurring", h.UpsertRecurring)
	r.Delete("/recurring/{id}", func(w http.ResponseWriter, r *http.Request) {
		h.DeleteRecurring(w, r, chi.URLParam(r, "id"))
	})
	r.Post("/sheet/close", h.CloseSheet)
}
