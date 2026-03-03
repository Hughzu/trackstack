package authwiring

import (
	"fmt"

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

	db, err := coredb.OpenLibSQL(dsn)
	if err != nil {
		return AuthDependencies{}, fmt.Errorf("auth db: %w", err)
	}

	store := authdb.NewSessionStore(db)
	service := auth.NewService(store, auth.Config{
		SessionIdleSeconds:          cfg.AuthSessionIdleSeconds,
		SessionAbsoluteSeconds:      cfg.AuthSessionAbsoluteSeconds,
		SessionRotateAfterSeconds:   cfg.AuthSessionRotateAfterSeconds,
		SessionRotationGraceSeconds: cfg.AuthSessionRotationGraceSeconds,
		SessionTouchSeconds:         cfg.AuthSessionTouchSeconds,
	})

	return AuthDependencies{
		Service: service,
		Close:   db.Close,
	}, nil
}
