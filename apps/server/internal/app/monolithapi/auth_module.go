package monolithapi

import (
	"context"
	"database/sql"
	"net/http"

	authhttpserver "github.com/Hughzu/trackstack/apps/server/internal/contexts/auth/adapters/inbound/http"
	authdb "github.com/Hughzu/trackstack/apps/server/internal/contexts/auth/adapters/outbound/db"
	"github.com/Hughzu/trackstack/apps/server/internal/contexts/auth/adapters/outbound/jwt"
	authtoken "github.com/Hughzu/trackstack/apps/server/internal/contexts/auth/adapters/outbound/token"
	authservices "github.com/Hughzu/trackstack/apps/server/internal/contexts/auth/application/services"
	usersservice "github.com/Hughzu/trackstack/apps/server/internal/contexts/users/application/services"
	"github.com/Hughzu/trackstack/apps/server/internal/platform/config"
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
	cfg config.Config,
) *AuthModule {
	userService := buildUsersModule(db)
	userProvider := &userProviderAdapter{svc: userService}
	tokenIssuer := jwt.NewIssuer(cfg.JWTSecret, cfg.AccessTokenTTL())
	sessionRepo := authdb.NewSessionRepository(db)
	refreshTokenManager := authtoken.NewManager()
	authService := authservices.NewAuthService(
		userProvider,
		tokenIssuer,
		sessionRepo,
		refreshTokenManager,
		cfg.RefreshTokenTTL(),
		cfg.RefreshTokenAbsoluteTTL(),
	)

	return &AuthModule{
		Handler: authhttpserver.NewAuthHandler(authService, authhttpserver.CookieConfig{
			Name:     cfg.RefreshCookieName,
			Domain:   cfg.RefreshCookieDomain,
			Path:     "/api/auth",
			Secure:   cfg.RefreshCookieSecure,
			SameSite: http.SameSiteLaxMode,
		}),
		Middleware: platformmiddleware.ResolveSession(cfg.JWTSecret),
	}
}
