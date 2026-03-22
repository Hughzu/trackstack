package bootstrap

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/Hughzu/trackstack/apps/server-next/internal/platform/config"
	"github.com/Hughzu/trackstack/apps/server-next/internal/platform/db"
)

type Databases struct {
	Calories *sql.DB
	Heat     *sql.DB
}

func connectDatabases(cfg config.Config) (*Databases, error) {
	poolCfg := db.PoolConfig{
		MaxOpenConns:    cfg.DBMaxOpenConns,
		MaxIdleConns:    cfg.DBMaxIdleConns,
		ConnMaxLifetime: time.Duration(cfg.DBConnMaxLifetimeSeconds) * time.Second,
		ConnMaxIdleTime: time.Duration(cfg.DBConnMaxIdleTimeSeconds) * time.Second,
	}

	caloriesDB, err := db.Open(cfg.TursoCaloriesURLHTTP, cfg.TursoCaloriesToken, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to open calories db: %w", err)
	}

	heatDB, err := db.Open(cfg.TursoHeatURLHTTP, cfg.TursoHeatToken, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to open heat db: %w", err)
	}

	return &Databases{
		Calories: caloriesDB,
		Heat:     heatDB,
	}, nil
}

func (d *Databases) CloseAll() []func() error {
	return []func() error{
		func() error { return d.Calories.Close() },
		func() error { return d.Heat.Close() },
	}
}
