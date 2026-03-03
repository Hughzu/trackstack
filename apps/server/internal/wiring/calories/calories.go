package calorieswiring

import (
	"fmt"

	"github.com/Hughzu/trackstack/apps/server/internal/core/config"
	coredb "github.com/Hughzu/trackstack/apps/server/internal/core/db"
	"github.com/Hughzu/trackstack/apps/server/internal/modules/calories"
	caloriesdb "github.com/Hughzu/trackstack/apps/server/internal/modules/calories/adapters/db"
	"github.com/Hughzu/trackstack/apps/server/internal/wiring/common"
)

type CaloriesDependencies struct {
	Service *calories.Service
	Close   func() error
}

func BuildCalories(cfg config.Config) (CaloriesDependencies, error) {
	dsn, err := common.BuildTursoDSN(common.TursoConfig{
		Mode:    cfg.DBConnectionMode,
		URL:     cfg.TursoHeatURL,
		URLHTTP: cfg.TursoHeatURLHTTP,
		URLWS:   cfg.TursoHeatURLWS,
		Token:   cfg.TursoHeatToken,
	})
	if err != nil {
		return CaloriesDependencies{}, fmt.Errorf("calories dsn: %w", err)
	}

	db, err := coredb.OpenLibSQL(dsn)
	if err != nil {
		return CaloriesDependencies{}, fmt.Errorf("calories db: %w", err)
	}

	store := caloriesdb.NewCaloriesStore(db)
	service := calories.NewService(store)

	return CaloriesDependencies{
		Service: service,
		Close:   db.Close,
	}, nil
}
