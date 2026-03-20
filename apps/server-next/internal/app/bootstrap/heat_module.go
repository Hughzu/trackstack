package bootstrap

import (
	"database/sql"

	heathttp "github.com/Hughzu/trackstack/apps/server-next/internal/contexts/heat/adapters/inbound/http"
	heatdb "github.com/Hughzu/trackstack/apps/server-next/internal/contexts/heat/adapters/outbound/db"
	heatservice "github.com/Hughzu/trackstack/apps/server-next/internal/contexts/heat/application/services"
)

func buildHeatModule(db *sql.DB) *heathttp.RefillHandler {

	// Outbound
	refillRepo := heatdb.NewRefillRepository(db)

	// Core
	useCase := heatservice.NewRefillService(refillRepo)

	// Inbound
	handler := heathttp.NewRefillHandler(useCase)

	return handler
}
