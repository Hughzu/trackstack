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

	return r
}
