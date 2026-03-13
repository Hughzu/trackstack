package calorieswiring

import (
	"fmt"
	"time"

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
		URL:     cfg.TursoCaloriesURL,
		URLHTTP: cfg.TursoCaloriesURLHTTP,
		URLWS:   cfg.TursoCaloriesURLWS,
		Token:   cfg.TursoCaloriesToken,
	})
	if err != nil {
		return CaloriesDependencies{}, fmt.Errorf("calories dsn: %w", err)
	}

	db, err := coredb.OpenLibSQL(dsn, coredb.PoolConfig{
		MaxOpenConns:    cfg.DBMaxOpenConns,
		MaxIdleConns:    cfg.DBMaxIdleConns,
		ConnMaxLifetime: time.Duration(cfg.DBConnMaxLifetimeSeconds) * time.Second,
		ConnMaxIdleTime: time.Duration(cfg.DBConnMaxIdleTimeSeconds) * time.Second,
	})
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
