package httptransport

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(handlers Handlers) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)

	r.Get("/health", Health)
	r.Get("/openapi.yaml", OpenAPISpec)

	r.Route("/api/heat", func(r chi.Router) {
		r.Get("/refills", handlers.Heat.ListRefills)
		r.Post("/refills", handlers.Heat.CreateRefill)
		r.Delete("/refills", handlers.Heat.DeleteRefill)
	})

	r.Route("/api/expenses", func(r chi.Router) {
		r.Get("/settings", handlers.Expenses.GetSettings)
		r.Post("/settings", handlers.Expenses.UpdateSettings)
		r.Get("/sheet/current", handlers.Expenses.GetCurrentSheet)
		r.Post("/entries", handlers.Expenses.AddExpense)
		r.Delete("/entries", handlers.Expenses.DeleteExpense)
		r.Post("/checklists", handlers.Expenses.UpsertChecklist)
		r.Delete("/checklists", handlers.Expenses.DeleteChecklist)
		r.Post("/checklists/complete", handlers.Expenses.CompleteChecklistItem)
		r.Post("/recurring", handlers.Expenses.UpsertRecurring)
		r.Delete("/recurring", handlers.Expenses.DeleteRecurring)
		r.Post("/sheet/close", handlers.Expenses.CloseSheet)
	})

	r.Route("/api/calories", func(r chi.Router) {
		r.Post("/log", handlers.Calories.AddLog)
		r.Delete("/log", handlers.Calories.DeleteLog)
		r.Get("/target", handlers.Calories.GetTarget)
		r.Post("/target", handlers.Calories.UpdateTarget)
	})

	r.Route("/api/auth", func(r chi.Router) {
		r.Post("/login", handlers.Auth.Login)
		r.Post("/logout", handlers.Auth.Logout)
	})

	return r
}
