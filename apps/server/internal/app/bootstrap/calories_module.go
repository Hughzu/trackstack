package bootstrap

import (
	"database/sql"

	calorieshttp "github.com/Hughzu/trackstack/apps/server/internal/contexts/calories/adapters/inbound/http"
	caloriesdb "github.com/Hughzu/trackstack/apps/server/internal/contexts/calories/adapters/outbound/db"
	caloriesservice "github.com/Hughzu/trackstack/apps/server/internal/contexts/calories/application/services"
)

func buildCaloriesModule(db *sql.DB) *calorieshttp.CaloriesHandler {
	// Outbound
	repository := caloriesdb.NewRepository(db)

	// Core
	targetUseCase := caloriesservice.NewTargetService(repository)
	logUseCase := caloriesservice.NewLogService(repository)
	dashboardUseCase := caloriesservice.NewDashboardService(repository, repository)

	// Inbound
	handler := calorieshttp.NewCaloriesHandler(targetUseCase, logUseCase, dashboardUseCase)

	return handler
}
