package userswiring

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/Hughzu/trackstack/apps/server/internal/core/config"
	coredb "github.com/Hughzu/trackstack/apps/server/internal/core/db"
	"github.com/Hughzu/trackstack/apps/server/internal/modules/users"
	usersdb "github.com/Hughzu/trackstack/apps/server/internal/modules/users/adapters/db"
	"github.com/Hughzu/trackstack/apps/server/internal/wiring/common"
)

type UsersDependencies struct {
	Service *users.Service
	Close   func() error
}

func OpenUsersDB(cfg config.Config) (*sql.DB, error) {
	dsn, err := common.BuildTursoDSN(common.TursoConfig{
		Mode:    cfg.DBConnectionMode,
		URL:     cfg.TursoUsersURL,
		URLHTTP: cfg.TursoUsersURLHTTP,
		URLWS:   cfg.TursoUsersURLWS,
		Token:   cfg.TursoUsersToken,
	})
	if err != nil {
		return nil, fmt.Errorf("users dsn: %w", err)
	}

	db, err := coredb.OpenLibSQL(dsn, coredb.PoolConfig{
		MaxOpenConns:    cfg.DBMaxOpenConns,
		MaxIdleConns:    cfg.DBMaxIdleConns,
		ConnMaxLifetime: time.Duration(cfg.DBConnMaxLifetimeSeconds) * time.Second,
		ConnMaxIdleTime: time.Duration(cfg.DBConnMaxIdleTimeSeconds) * time.Second,
	})
	if err != nil {
		return nil, fmt.Errorf("users db: %w", err)
	}

	return db, nil
}

func NewUsersService(db *sql.DB) *users.Service {
	store := usersdb.NewUsersStore(db)
	return users.NewService(store)
}

func BuildUsers(cfg config.Config) (UsersDependencies, error) {
	db, err := OpenUsersDB(cfg)
	if err != nil {
		return UsersDependencies{}, err
	}

	service := NewUsersService(db)

	return UsersDependencies{
		Service: service,
		Close:   db.Close,
	}, nil
}
