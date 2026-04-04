package monolithapi

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Hughzu/trackstack/apps/server/internal/platform/config"
	"github.com/Hughzu/trackstack/apps/server/internal/platform/logging"
)

type Runtime struct {
	Config  config.Config
	Logger  *slog.Logger
	Handler http.Handler
	closers []func() error
}

func NewRuntime() (*Runtime, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("Error while loading config : %w", err)
	}

	logger := logging.New(cfg.LogLevel)
	slog.SetDefault(logger)

	dbs, err := connectDatabases(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to connect databases: %w", err)
	}

	heatHandler := buildHeatModule(dbs.Heat)
	caloriesHandler := buildCaloriesModule(dbs.Calories)
	expensesHandler := buildExpensesModule(dbs.Expenses)
	authModule := buildAuthModule(dbs.Users, cfg.JWTSecret)

	runtime := &Runtime{
		Config:  cfg,
		Logger:  logger,
		Handler: newRouter(logger, cfg.CORSAllowedOrigin, authModule, heatHandler, caloriesHandler, expensesHandler),
		closers: dbs.CloseAll(),
	}

	return runtime, nil
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
