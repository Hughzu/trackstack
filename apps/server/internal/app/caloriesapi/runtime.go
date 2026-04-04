package caloriesapi

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	calorieshttp "github.com/Hughzu/trackstack/apps/server/internal/contexts/calories/adapters/inbound/http"
	caloriesdb "github.com/Hughzu/trackstack/apps/server/internal/contexts/calories/adapters/outbound/db"
	caloriesservice "github.com/Hughzu/trackstack/apps/server/internal/contexts/calories/application/services"
	platformdb "github.com/Hughzu/trackstack/apps/server/internal/platform/db"
	"github.com/Hughzu/trackstack/apps/server/internal/platform/logging"
	"github.com/Hughzu/trackstack/apps/server/internal/platform/middleware"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

type Config struct {
	Env               string
	Port              string
	LogLevel          string
	CORSAllowedOrigin string
	JWTSecret         string

	TursoCaloriesURLHTTP string
	TursoCaloriesToken   string

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

func NewRuntime() (*Runtime, error) {
	cfg, err := LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("load calories runtime config: %w", err)
	}

	logger := logging.New(cfg.LogLevel)
	slog.SetDefault(logger)

	caloriesDB, err := connectDatabase(cfg)
	if err != nil {
		return nil, fmt.Errorf("connect calories database: %w", err)
	}

	caloriesHandler := buildCaloriesModule(caloriesDB)

	return &Runtime{
		Config:  cfg,
		Logger:  logger,
		Handler: newRouter(logger, cfg.CORSAllowedOrigin, cfg.JWTSecret, caloriesHandler),
		closers: []func() error{caloriesDB.Close},
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
		Env:                  getEnv("APP_ENV", "local"),
		Port:                 getEnv("PORT", "8080"),
		LogLevel:             getEnv("LOG_LEVEL", "info"),
		CORSAllowedOrigin:    getEnv("CORS_ALLOWED_ORIGIN", ""),
		JWTSecret:            getEnv("JWT_SECRET", ""),
		TursoCaloriesURLHTTP: getEnv("TURSO_CALORIES_URL_HTTP", ""),
		TursoCaloriesToken:   getEnv("TURSO_CALORIES_TOKEN", ""),
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

	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func (cfg Config) Validate() error {
	if strings.TrimSpace(cfg.Port) == "" {
		return fmt.Errorf("PORT must not be empty")
	}
	if strings.TrimSpace(cfg.TursoCaloriesURLHTTP) == "" {
		return fmt.Errorf("TURSO_CALORIES_URL_HTTP must not be empty")
	}
	if strings.TrimSpace(cfg.TursoCaloriesToken) == "" {
		return fmt.Errorf("TURSO_CALORIES_TOKEN must not be empty")
	}
	if strings.TrimSpace(cfg.JWTSecret) == "" {
		return fmt.Errorf("JWT_SECRET must not be empty")
	}

	return nil
}

func connectDatabase(cfg Config) (*sql.DB, error) {
	poolCfg := platformdb.PoolConfig{
		MaxOpenConns:    cfg.DBMaxOpenConns,
		MaxIdleConns:    cfg.DBMaxIdleConns,
		ConnMaxLifetime: time.Duration(cfg.DBConnMaxLifetimeSeconds) * time.Second,
		ConnMaxIdleTime: time.Duration(cfg.DBConnMaxIdleTimeSeconds) * time.Second,
	}

	return platformdb.Open(cfg.TursoCaloriesURLHTTP, cfg.TursoCaloriesToken, poolCfg)
}

func buildCaloriesModule(db *sql.DB) *calorieshttp.CaloriesHandler {
	repository := caloriesdb.NewRepository(db)
	targetUseCase := caloriesservice.NewTargetService(repository)
	logUseCase := caloriesservice.NewLogService(repository)
	dashboardUseCase := caloriesservice.NewDashboardService(repository, repository)
	return calorieshttp.NewCaloriesHandler(targetUseCase, logUseCase, dashboardUseCase)
}

func newRouter(logger *slog.Logger, corsAllowedOrigin string, jwtSecret string, caloriesHandler *calorieshttp.CaloriesHandler) http.Handler {
	r := chi.NewRouter()

	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(chimiddleware.Recoverer)
	r.Use(middleware.Logger(logger))
	r.Use(middleware.CORS(corsAllowedOrigin))

	r.Get("/health", health)
	r.Route("/api", func(r chi.Router) {
		r.Get("/health", health)
		r.Group(func(r chi.Router) {
			r.Use(middleware.ResolveSession(jwtSecret))
			r.Route("/calories", caloriesHandler.RegisterRoutes)
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
