package monolithapi

import (
	"database/sql"

	usersdb "github.com/Hughzu/trackstack/apps/server/internal/contexts/users/adapters/outbound/db"
	usersservice "github.com/Hughzu/trackstack/apps/server/internal/contexts/users/application/services"
)

func buildUsersModule(db *sql.DB) *usersservice.UserService {
	// Outbound
	usersRepo := usersdb.NewUserRepository(db)

	// Core
	userService := usersservice.NewUserService(usersRepo)
	return userService
}
