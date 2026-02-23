package heatwiring

import (
	"fmt"

	"github.com/Hughzu/trackstack/apps/server/internal/core/config"
	coredb "github.com/Hughzu/trackstack/apps/server/internal/core/db"
	"github.com/Hughzu/trackstack/apps/server/internal/modules/heat"
	heatdb "github.com/Hughzu/trackstack/apps/server/internal/modules/heat/adapters/db"
	"github.com/Hughzu/trackstack/apps/server/internal/wiring/common"
)

type HeatDependencies struct {
	Service *heat.Service
	Close   func() error
}

func BuildHeat(cfg config.Config) (HeatDependencies, error) {
	refillDSN, err := common.BuildTursoDSN(common.TursoConfig{
		Mode:    cfg.DBConnectionMode,
		URL:     cfg.TursoHeatURL,
		URLHTTP: cfg.TursoHeatURLHTTP,
		URLWS:   cfg.TursoHeatURLWS,
		Token:   cfg.TursoHeatToken,
	})
	if err != nil {
		return HeatDependencies{}, fmt.Errorf("heat dsn: %w", err)
	}

	heatDB, err := coredb.OpenLibSQL(refillDSN)
	if err != nil {
		return HeatDependencies{}, fmt.Errorf("heat db: %w", err)
	}

	store := heatdb.NewRefillStore(heatDB)
	service := heat.NewService(store)

	return HeatDependencies{
		Service: service,
		Close:   heatDB.Close,
	}, nil
}
