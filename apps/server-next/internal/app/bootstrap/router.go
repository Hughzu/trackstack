package bootstrap

import (
	"log/slog"
	"net/http"

	caloriesinboundhttp "github.com/Hughzu/trackstack/apps/server-next/internal/contexts/calories/adapters/inbound/http"
	heatinboundhttp "github.com/Hughzu/trackstack/apps/server-next/internal/contexts/heat/adapters/inbound/http"
	"github.com/Hughzu/trackstack/apps/server-next/internal/platform/middleware"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

func newRouter(logger *slog.Logger, heatHandler *heatinboundhttp.RefillHandler, caloriesHandler *caloriesinboundhttp.CaloriesHandler) http.Handler {
	r := chi.NewRouter()

	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(chimiddleware.Recoverer)
	r.Use(middleware.Logger(logger))

	r.Get("/health", health)
	r.Route("/api", func(r chi.Router) {
		r.Get("/health", health)

		r.Route("/calories", caloriesHandler.RegisterRoutes)
		r.Route("/heat", heatHandler.RegisterRoutes)
	})
	return r
}
func health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, _ = w.Write([]byte(`{"status":"ok"}`))
}
