package app

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Hughzu/trackstack/apps/server/internal/core/config"
	"github.com/Hughzu/trackstack/apps/server/internal/core/logging"
	httptransport "github.com/Hughzu/trackstack/apps/server/internal/transport/http"
	authwiring "github.com/Hughzu/trackstack/apps/server/internal/wiring/auth"
	calorieswiring "github.com/Hughzu/trackstack/apps/server/internal/wiring/calories"
	commonwiring "github.com/Hughzu/trackstack/apps/server/internal/wiring/common"
	expenseswiring "github.com/Hughzu/trackstack/apps/server/internal/wiring/expenses"
	heatwiring "github.com/Hughzu/trackstack/apps/server/internal/wiring/heat"
	userswiring "github.com/Hughzu/trackstack/apps/server/internal/wiring/users"
)

type Runtime struct {
	Config  config.Config
	Logger  *slog.Logger
	Handler http.Handler
	closers []resourceCloser
}

type resourceCloser struct {
	name  string
	close func() error
}

func NewRuntime() (*Runtime, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}

	logger := logging.New(cfg.LogLevel)
	slog.SetDefault(logger)

	logTursoSelection(logger, "calories", commonwiring.TursoConfig{Mode: cfg.DBConnectionMode, URL: cfg.TursoCaloriesURL, URLHTTP: cfg.TursoCaloriesURLHTTP, URLWS: cfg.TursoCaloriesURLWS})
	logTursoSelection(logger, "expenses", commonwiring.TursoConfig{Mode: cfg.DBConnectionMode, URL: cfg.TursoExpensesURL, URLHTTP: cfg.TursoExpensesURLHTTP, URLWS: cfg.TursoExpensesURLWS})
	logTursoSelection(logger, "heat", commonwiring.TursoConfig{Mode: cfg.DBConnectionMode, URL: cfg.TursoHeatURL, URLHTTP: cfg.TursoHeatURLHTTP, URLWS: cfg.TursoHeatURLWS})
	logTursoSelection(logger, "users", commonwiring.TursoConfig{Mode: cfg.DBConnectionMode, URL: cfg.TursoUsersURL, URLHTTP: cfg.TursoUsersURLHTTP, URLWS: cfg.TursoUsersURLWS})

	runtime := &Runtime{
		Config:  cfg,
		Logger:  logger,
		closers: make([]resourceCloser, 0, 4),
	}

	heatDeps, err := heatwiring.BuildHeat(cfg)
	if err != nil {
		return nil, fmt.Errorf("build heat dependencies: %w", err)
	}
	runtime.addCloser("heat db", heatDeps.Close)

	expensesDeps, err := expenseswiring.BuildExpenses(cfg)
	if err != nil {
		_ = runtime.Close()
		return nil, fmt.Errorf("build expenses dependencies: %w", err)
	}
	runtime.addCloser("expenses db", expensesDeps.Close)

	caloriesDeps, err := calorieswiring.BuildCalories(cfg)
	if err != nil {
		_ = runtime.Close()
		return nil, fmt.Errorf("build calories dependencies: %w", err)
	}
	runtime.addCloser("calories db", caloriesDeps.Close)

	usersDB, err := userswiring.OpenUsersDB(cfg)
	if err != nil {
		_ = runtime.Close()
		return nil, fmt.Errorf("open users/auth db: %w", err)
	}
	runtime.addCloser("users/auth db", usersDB.Close)

	usersService := userswiring.NewUsersService(usersDB)
	authService := authwiring.NewAuthService(usersDB, cfg)

	handlers := httptransport.NewHandlers(httptransport.Deps{
		HeatService:        heatDeps.Service,
		ExpensesService:    expensesDeps.Service,
		CaloriesService:    caloriesDeps.Service,
		UsersService:       usersService,
		AuthService:        authService,
		AuthCookieName:     cfg.AuthCookieName,
		AuthCookieSecure:   cfg.AuthCookieSecure,
		AuthCookieSameSite: cfg.AuthCookieSameSite,
	})

	runtime.Handler = httptransport.NewRouter(handlers, cfg.CORSAllowedOrigin)

	return runtime, nil
}

func (r *Runtime) Close() error {
	if r == nil {
		return nil
	}

	var errs []error
	for i := len(r.closers) - 1; i >= 0; i-- {
		closer := r.closers[i]
		if err := closer.close(); err != nil {
			err = fmt.Errorf("%s: %w", closer.name, err)
			errs = append(errs, err)
			if r.Logger != nil {
				r.Logger.Error("resource close error", "resource", closer.name, "error", err)
			}
		}
	}

	r.closers = nil

	return errors.Join(errs...)
}

func (r *Runtime) addCloser(name string, closeFn func() error) {
	r.closers = append(r.closers, resourceCloser{name: name, close: closeFn})
}

func logTursoSelection(logger *slog.Logger, domain string, cfg commonwiring.TursoConfig) {
	selection, err := commonwiring.ResolveTursoSelection(cfg)
	if err != nil {
		logger.Error("turso config error", "domain", domain, "error", err)
		return
	}

	logger.Info("turso connection selected", "domain", domain, "mode", selection.Mode, "scheme", selection.Scheme, "host", selection.Host)
}
