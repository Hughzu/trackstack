package bootstrap

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/Hughzu/trackstack/apps/server/internal/platform/config"
	"github.com/Hughzu/trackstack/apps/server/internal/platform/db"
)

type Databases struct {
	Calories *sql.DB
	Expenses *sql.DB
	Heat     *sql.DB
	Users    *sql.DB
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

	expensesDB, err := db.Open(cfg.TursoExpensesURLHTTP, cfg.TursoExpensesToken, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to open expenses db: %w", err)
	}

	usersDB, err := db.Open(cfg.TursoUsersURLHTTP, cfg.TursoUsersToken, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to open users db: %w", err)
	}

	return &Databases{
		Calories: caloriesDB,
		Expenses: expensesDB,
		Heat:     heatDB,
		Users:    usersDB,
	}, nil
}

func (d *Databases) CloseAll() []func() error {
	return []func() error{
		func() error { return d.Calories.Close() },
		func() error { return d.Expenses.Close() },
		func() error { return d.Heat.Close() },
		func() error { return d.Users.Close() },
	}
}
