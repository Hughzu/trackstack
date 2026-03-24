package bootstrap

import (
	"context"
	"database/sql"
	"net/http"

	authhttpserver "github.com/Hughzu/trackstack/apps/server/internal/contexts/auth/adapters/inbound/http"
	"github.com/Hughzu/trackstack/apps/server/internal/contexts/auth/adapters/outbound/jwt"
	authservices "github.com/Hughzu/trackstack/apps/server/internal/contexts/auth/application/services"
	usersservice "github.com/Hughzu/trackstack/apps/server/internal/contexts/users/application/services"
	platformmiddleware "github.com/Hughzu/trackstack/apps/server/internal/platform/middleware"
)

type AuthModule struct {
	Handler    *authhttpserver.AuthHandler
	Middleware func(http.Handler) http.Handler
}

type userProviderAdapter struct {
	svc *usersservice.UserService
}

func (a *userProviderAdapter) VerifyCredentials(ctx context.Context, email string, password string) (string, error) {
	u, err := a.svc.VerifyCredentials(ctx, email, password)
	return u.ID, err
}

func (a *userProviderAdapter) UpdateLastLogin(ctx context.Context, userID string, timestamp string) error {
	return a.svc.UpdateLastLogin(ctx, userID, timestamp)
}

func buildAuthModule(
	db *sql.DB,
	jwtSecret string,
) *AuthModule {
	userService := buildUsersModule(db)
	userProvider := &userProviderAdapter{svc: userService}
	tokenIssuer := jwt.NewIssuer(jwtSecret)
	authService := authservices.NewAuthService(userProvider, tokenIssuer)

	return &AuthModule{
		Handler:    authhttpserver.NewAuthHandler(authService),
		Middleware: platformmiddleware.ResolveSession(jwtSecret),
	}
}
