package bootstrap

import (
	"net/http"

	heatinboundhttp "github.com/Hughzu/trackstack/apps/server-next/internal/contexts/heat/adapters/inbound/http"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func newRouter(heatHandler *heatinboundhttp.RefillHandler) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)

	r.Get("/health", health)
	r.Route("/api", func(r chi.Router) {
		r.Get("/health", health)

		r.Route("/heat", heatHandler.RegisterRoutes)
	})
	return r
}
func health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, _ = w.Write([]byte(`{"status":"ok"}`))
}
