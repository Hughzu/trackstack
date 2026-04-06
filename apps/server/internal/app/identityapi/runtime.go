package identityapi

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	authhttpserver "github.com/Hughzu/trackstack/apps/server/internal/contexts/auth/adapters/inbound/http"
	authdb "github.com/Hughzu/trackstack/apps/server/internal/contexts/auth/adapters/outbound/db"
	"github.com/Hughzu/trackstack/apps/server/internal/contexts/auth/adapters/outbound/jwt"
	authtoken "github.com/Hughzu/trackstack/apps/server/internal/contexts/auth/adapters/outbound/token"
	authservices "github.com/Hughzu/trackstack/apps/server/internal/contexts/auth/application/services"
	usersdb "github.com/Hughzu/trackstack/apps/server/internal/contexts/users/adapters/outbound/db"
	usersservice "github.com/Hughzu/trackstack/apps/server/internal/contexts/users/application/services"
	platformdb "github.com/Hughzu/trackstack/apps/server/internal/platform/db"
	"github.com/Hughzu/trackstack/apps/server/internal/platform/logging"
	"github.com/Hughzu/trackstack/apps/server/internal/platform/middleware"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

type Config struct {
	Env                          string
	Port                         string
	LogLevel                     string
	CORSAllowedOrigin            string
	JWTSecret                    string
	AccessTokenTTLMinutes        int
	RefreshTokenTTLHours         int
	RefreshTokenAbsoluteTTLHours int
	RefreshCookieName            string
	RefreshCookieSecure          bool
	RefreshCookieDomain          string

	TursoUsersURLHTTP string
	TursoUsersToken   string

	DBMaxOpenConns           int
	DBMaxIdleConns           int
	DBConnMaxLifetimeSeconds int
	DBConnMaxIdleTimeSeconds int
}

type Runtime struct {
	Config  Config
	Logger  *slog.Logger
	Handler http.Handler
	closers []func() error
}

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

func NewRuntime() (*Runtime, error) {
	cfg, err := LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("load identity runtime config: %w", err)
	}

	logger := logging.New(cfg.LogLevel)
	slog.SetDefault(logger)

	usersDB, err := connectDatabase(cfg)
	if err != nil {
		return nil, fmt.Errorf("connect users database: %w", err)
	}

	authModule := buildAuthModule(usersDB, cfg)

	return &Runtime{
		Config:  cfg,
		Logger:  logger,
		Handler: newRouter(logger, cfg.CORSAllowedOrigin, authModule),
		closers: []func() error{usersDB.Close},
	}, nil
}

func (r *Runtime) Close() error {
	if r == nil {
		return nil
	}

	var firstErr error
	for i := len(r.closers) - 1; i >= 0; i-- {
		if err := r.closers[i](); err != nil && firstErr == nil {
			firstErr = err
		}
	}

	r.closers = nil
	return firstErr
}

func LoadConfig() (Config, error) {
	var err error
	cfg := Config{
		Env:                 getEnv("APP_ENV", "local"),
		Port:                getEnv("PORT", "8080"),
		LogLevel:            getEnv("LOG_LEVEL", "info"),
		CORSAllowedOrigin:   getEnv("CORS_ALLOWED_ORIGIN", ""),
		JWTSecret:           getEnv("JWT_SECRET", ""),
		RefreshCookieName:   getEnv("REFRESH_COOKIE_NAME", "trackstack_refresh"),
		RefreshCookieDomain: getEnv("REFRESH_COOKIE_DOMAIN", ""),
		TursoUsersURLHTTP:   getEnv("TURSO_USERS_URL_HTTP", ""),
		TursoUsersToken:     getEnv("TURSO_USERS_TOKEN", ""),
	}

	cfg.DBMaxOpenConns, err = getEnvInt("DB_MAX_OPEN_CONNS", 10)
	if err != nil {
		return Config{}, err
	}

	cfg.DBMaxIdleConns, err = getEnvInt("DB_MAX_IDLE_CONNS", 5)
	if err != nil {
		return Config{}, err
	}

	cfg.DBConnMaxLifetimeSeconds, err = getEnvInt("DB_CONN_MAX_LIFETIME_SECONDS", 300)
	if err != nil {
		return Config{}, err
	}

	cfg.DBConnMaxIdleTimeSeconds, err = getEnvInt("DB_CONN_MAX_IDLE_TIME_SECONDS", 60)
	if err != nil {
		return Config{}, err
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

	cfg.RefreshCookieSecure = getEnvBool("REFRESH_COOKIE_SECURE", cfg.Env != "local")

	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func (cfg Config) Validate() error {
	if strings.TrimSpace(cfg.Port) == "" {
		return fmt.Errorf("PORT must not be empty")
	}
	if strings.TrimSpace(cfg.TursoUsersURLHTTP) == "" {
		return fmt.Errorf("TURSO_USERS_URL_HTTP must not be empty")
	}
	if strings.TrimSpace(cfg.TursoUsersToken) == "" {
		return fmt.Errorf("TURSO_USERS_TOKEN must not be empty")
	}
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

func connectDatabase(cfg Config) (*sql.DB, error) {
	poolCfg := platformdb.PoolConfig{
		MaxOpenConns:    cfg.DBMaxOpenConns,
		MaxIdleConns:    cfg.DBMaxIdleConns,
		ConnMaxLifetime: time.Duration(cfg.DBConnMaxLifetimeSeconds) * time.Second,
		ConnMaxIdleTime: time.Duration(cfg.DBConnMaxIdleTimeSeconds) * time.Second,
	}

	return platformdb.Open(cfg.TursoUsersURLHTTP, cfg.TursoUsersToken, poolCfg)
}

func buildAuthModule(db *sql.DB, cfg Config) *AuthModule {
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
		Middleware: middleware.ResolveSession(cfg.JWTSecret),
	}
}

func buildUsersModule(db *sql.DB) *usersservice.UserService {
	usersRepo := usersdb.NewUserRepository(db)
	return usersservice.NewUserService(usersRepo)
}

func newRouter(logger *slog.Logger, corsAllowedOrigin string, authModule *AuthModule) http.Handler {
	r := chi.NewRouter()

	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(chimiddleware.Recoverer)
	r.Use(middleware.Logger(logger))
	r.Use(middleware.CORS(corsAllowedOrigin))

	r.Get("/health", health)
	r.Route("/api", func(r chi.Router) {
		r.Get("/health", health)
		r.Route("/auth", func(r chi.Router) {
			authModule.Handler.RegisterRoutes(r, authModule.Middleware)
		})
	})

	return r
}

func health(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}

func getEnv(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}

	return fallback
}

func getEnvInt(key string, fallback int) (int, error) {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback, nil
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("invalid %s: %w", key, err)
	}

	return parsed, nil
}

func getEnvBool(key string, fallback bool) bool {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}

	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return fallback
	}

	return parsed
}
