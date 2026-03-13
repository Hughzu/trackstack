package authwiring

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/Hughzu/trackstack/apps/server/internal/core/config"
	coredb "github.com/Hughzu/trackstack/apps/server/internal/core/db"
	"github.com/Hughzu/trackstack/apps/server/internal/modules/auth"
	authdb "github.com/Hughzu/trackstack/apps/server/internal/modules/auth/adapters/db"
	"github.com/Hughzu/trackstack/apps/server/internal/wiring/common"
)

type AuthDependencies struct {
	Service *auth.Service
	Close   func() error
}

func NewAuthService(db *sql.DB, cfg config.Config) *auth.Service {
	store := authdb.NewSessionStore(db)
	return auth.NewService(store, auth.Config{
		SessionIdleSeconds:          cfg.AuthSessionIdleSeconds,
		SessionAbsoluteSeconds:      cfg.AuthSessionAbsoluteSeconds,
		SessionRotateAfterSeconds:   cfg.AuthSessionRotateAfterSeconds,
		SessionRotationGraceSeconds: cfg.AuthSessionRotationGraceSeconds,
		SessionTouchSeconds:         cfg.AuthSessionTouchSeconds,
	})
}

func BuildAuth(cfg config.Config) (AuthDependencies, error) {
	dsn, err := common.BuildTursoDSN(common.TursoConfig{
		Mode:    cfg.DBConnectionMode,
		URL:     cfg.TursoUsersURL,
		URLHTTP: cfg.TursoUsersURLHTTP,
		URLWS:   cfg.TursoUsersURLWS,
		Token:   cfg.TursoUsersToken,
	})
	if err != nil {
		return AuthDependencies{}, fmt.Errorf("auth dsn: %w", err)
	}

	db, err := coredb.OpenLibSQL(dsn, coredb.PoolConfig{
		MaxOpenConns:    cfg.DBMaxOpenConns,
		MaxIdleConns:    cfg.DBMaxIdleConns,
		ConnMaxLifetime: time.Duration(cfg.DBConnMaxLifetimeSeconds) * time.Second,
		ConnMaxIdleTime: time.Duration(cfg.DBConnMaxIdleTimeSeconds) * time.Second,
	})
	if err != nil {
		return AuthDependencies{}, fmt.Errorf("auth db: %w", err)
	}

	service := NewAuthService(db, cfg)

	return AuthDependencies{
		Service: service,
		Close:   db.Close,
	}, nil
}
