package authruntime

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"time"

	authhttpserver "github.com/Hughzu/trackstack/apps/server/internal/contexts/auth/adapters/inbound/http"
	authdb "github.com/Hughzu/trackstack/apps/server/internal/contexts/auth/adapters/outbound/db"
	"github.com/Hughzu/trackstack/apps/server/internal/contexts/auth/adapters/outbound/jwt"
	authtoken "github.com/Hughzu/trackstack/apps/server/internal/contexts/auth/adapters/outbound/token"
	authservices "github.com/Hughzu/trackstack/apps/server/internal/contexts/auth/application/services"
	usersdb "github.com/Hughzu/trackstack/apps/server/internal/contexts/users/adapters/outbound/db"
	usersservice "github.com/Hughzu/trackstack/apps/server/internal/contexts/users/application/services"
	platformconfig "github.com/Hughzu/trackstack/apps/server/internal/platform/config"
	platformmiddleware "github.com/Hughzu/trackstack/apps/server/internal/platform/middleware"
)

type Config struct {
	JWTSecret                    string
	AccessTokenTTLMinutes        int
	RefreshTokenTTLHours         int
	RefreshTokenAbsoluteTTLHours int
	RefreshCookieName            string
	RefreshCookieSecure          bool
	RefreshCookieDomain          string
}

type Module struct {
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

func NewModule(db *sql.DB, cfg Config) *Module {
	userService := usersservice.NewUserService(usersdb.NewUserRepository(db))
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

	return &Module{
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

func LoadConfig(getEnv func(string, string) string, getEnvInt func(string, int) (int, error), getEnvBool func(string, bool) bool, appEnv string) (Config, error) {
	var err error
	cfg := Config{
		JWTSecret:           getEnv("JWT_SECRET", ""),
		RefreshCookieName:   getEnv("REFRESH_COOKIE_NAME", "trackstack_refresh"),
		RefreshCookieDomain: getEnv("REFRESH_COOKIE_DOMAIN", ""),
	}

	cfg.AccessTokenTTLMinutes, err = getEnvInt("ACCESS_TOKEN_TTL_MINUTES", 15)
	if err != nil {
		return Config{}, err
	}

	cfg.RefreshTokenTTLHours, err = getEnvInt("REFRESH_TOKEN_TTL_HOURS", 24*30)
	if err != nil {
		return Config{}, err
	}

	cfg.RefreshTokenAbsoluteTTLHours, err = getEnvInt("REFRESH_TOKEN_ABSOLUTE_TTL_HOURS", 24*30)
	if err != nil {
		return Config{}, err
	}

	cfg.RefreshCookieSecure = getEnvBool("REFRESH_COOKIE_SECURE", appEnv != "local")

	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func FromPlatformConfig(cfg platformconfig.Config) Config {
	return Config{
		JWTSecret:                    cfg.JWTSecret,
		AccessTokenTTLMinutes:        cfg.AccessTokenTTLMinutes,
		RefreshTokenTTLHours:         cfg.RefreshTokenTTLHours,
		RefreshTokenAbsoluteTTLHours: cfg.RefreshTokenAbsoluteTTLHours,
		RefreshCookieName:            cfg.RefreshCookieName,
		RefreshCookieSecure:          cfg.RefreshCookieSecure,
		RefreshCookieDomain:          cfg.RefreshCookieDomain,
	}
}

func (cfg Config) Validate() error {
	if strings.TrimSpace(cfg.JWTSecret) == "" {
		return fmt.Errorf("JWT_SECRET must not be empty")
	}
	if cfg.AccessTokenTTLMinutes <= 0 {
		return fmt.Errorf("ACCESS_TOKEN_TTL_MINUTES must be greater than zero")
	}
	if cfg.RefreshTokenTTLHours <= 0 {
		return fmt.Errorf("REFRESH_TOKEN_TTL_HOURS must be greater than zero")
	}
	if cfg.RefreshTokenAbsoluteTTLHours < cfg.RefreshTokenTTLHours {
		return fmt.Errorf("REFRESH_TOKEN_ABSOLUTE_TTL_HOURS must be greater than or equal to REFRESH_TOKEN_TTL_HOURS")
	}
	if strings.TrimSpace(cfg.RefreshCookieName) == "" {
		return fmt.Errorf("REFRESH_COOKIE_NAME must not be empty")
	}

	return nil
}

func (cfg Config) AccessTokenTTL() time.Duration {
	return time.Duration(cfg.AccessTokenTTLMinutes) * time.Minute
}

func (cfg Config) RefreshTokenTTL() time.Duration {
	return time.Duration(cfg.RefreshTokenTTLHours) * time.Hour
}

func (cfg Config) RefreshTokenAbsoluteTTL() time.Duration {
	return time.Duration(cfg.RefreshTokenAbsoluteTTLHours) * time.Hour
}
