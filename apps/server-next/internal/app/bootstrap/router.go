package bootstrap

import (
	"log/slog"
	"net/http"

	caloriesinboundhttp "github.com/Hughzu/trackstack/apps/server-next/internal/contexts/calories/adapters/inbound/http"
	expensesinboundhttp "github.com/Hughzu/trackstack/apps/server-next/internal/contexts/expenses/adapters/inbound/http"
	heatinboundhttp "github.com/Hughzu/trackstack/apps/server-next/internal/contexts/heat/adapters/inbound/http"
	"github.com/Hughzu/trackstack/apps/server-next/internal/platform/middleware"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

func newRouter(
	logger *slog.Logger,
	corsAllowedOrigin string,
	authModule *AuthModule,
	heatHandler *heatinboundhttp.RefillHandler,
	caloriesHandler *caloriesinboundhttp.CaloriesHandler,
	expensesHandler *expensesinboundhttp.ExpensesHandler,
) http.Handler {
	r := chi.NewRouter()

	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(chimiddleware.Recoverer)
	r.Use(middleware.Logger(logger))
	r.Use(middleware.CORS(corsAllowedOrigin))

	r.Get("/health", health)
	r.Get("/openapi.yaml", openapi)
	r.Route("/api", func(r chi.Router) {
		r.Get("/health", health)

		r.Route("/auth", func(r chi.Router) {
			authModule.Handler.RegisterRoutes(r, authModule.Middleware)
		})

		r.Group(func(r chi.Router) {
			r.Use(authModule.Middleware)

			r.Route("/calories", caloriesHandler.RegisterRoutes)
			r.Route("/expenses", expensesHandler.RegisterRoutes)
			r.Route("/heat", heatHandler.RegisterRoutes)
		})
	})
	return r
}
func health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, _ = w.Write([]byte(`{"status":"ok"}`))
}
