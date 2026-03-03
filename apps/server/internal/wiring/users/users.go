package userswiring

import (
	"fmt"

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

func BuildUsers(cfg config.Config) (UsersDependencies, error) {
	dsn, err := common.BuildTursoDSN(common.TursoConfig{
		Mode:    cfg.DBConnectionMode,
		URL:     cfg.TursoUsersURL,
		URLHTTP: cfg.TursoUsersURLHTTP,
		URLWS:   cfg.TursoUsersURLWS,
		Token:   cfg.TursoUsersToken,
	})
	if err != nil {
		return UsersDependencies{}, fmt.Errorf("users dsn: %w", err)
	}

	db, err := coredb.OpenLibSQL(dsn)
	if err != nil {
		return UsersDependencies{}, fmt.Errorf("users db: %w", err)
	}

	store := usersdb.NewUsersStore(db)
	service := users.NewService(store)

	return UsersDependencies{
		Service: service,
		Close:   db.Close,
	}, nil
}
