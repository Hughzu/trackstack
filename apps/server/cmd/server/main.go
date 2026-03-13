package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Hughzu/trackstack/apps/server/internal/core/config"
	"github.com/Hughzu/trackstack/apps/server/internal/core/logging"
	httptransport "github.com/Hughzu/trackstack/apps/server/internal/transport/http"
	authwiring "github.com/Hughzu/trackstack/apps/server/internal/wiring/auth"
	calorieswiring "github.com/Hughzu/trackstack/apps/server/internal/wiring/calories"
	commonwiring "github.com/Hughzu/trackstack/apps/server/internal/wiring/common"
	expenseswiring "github.com/Hughzu/trackstack/apps/server/internal/wiring/expenses"
	heatwiring "github.com/Hughzu/trackstack/apps/server/internal/wiring/heat"
	userswiring "github.com/Hughzu/trackstack/apps/server/internal/wiring/users"

	"github.com/joho/godotenv"
)

func init() {
	_ = godotenv.Load()
}

func main() {
	cfg, err := config.Load()
	if err != nil {
		_, _ = os.Stderr.WriteString("config error: " + err.Error() + "\n")
		os.Exit(1)
	}

	logger := logging.New(cfg.LogLevel)
	slog.SetDefault(logger)
	logTursoSelection(logger, "calories", commonwiring.TursoConfig{Mode: cfg.DBConnectionMode, URL: cfg.TursoCaloriesURL, URLHTTP: cfg.TursoCaloriesURLHTTP, URLWS: cfg.TursoCaloriesURLWS})
	logTursoSelection(logger, "expenses", commonwiring.TursoConfig{Mode: cfg.DBConnectionMode, URL: cfg.TursoExpensesURL, URLHTTP: cfg.TursoExpensesURLHTTP, URLWS: cfg.TursoExpensesURLWS})
	logTursoSelection(logger, "heat", commonwiring.TursoConfig{Mode: cfg.DBConnectionMode, URL: cfg.TursoHeatURL, URLHTTP: cfg.TursoHeatURLHTTP, URLWS: cfg.TursoHeatURLWS})
	logTursoSelection(logger, "users", commonwiring.TursoConfig{Mode: cfg.DBConnectionMode, URL: cfg.TursoUsersURL, URLHTTP: cfg.TursoUsersURLHTTP, URLWS: cfg.TursoUsersURLWS})

	heatDeps, err := heatwiring.BuildHeat(cfg)
	if err != nil {
		logger.Error("heat db error", "error", err)
		os.Exit(1)
	}
	defer func() {
		if err := heatDeps.Close(); err != nil {
			logger.Error("heat db close error", "error", err)
		}
	}()

	expensesDeps, err := expenseswiring.BuildExpenses(cfg)
	if err != nil {
		logger.Error("expenses db error", "error", err)
		os.Exit(1)
	}
	defer func() {
		if err := expensesDeps.Close(); err != nil {
			logger.Error("expenses db close error", "error", err)
		}
	}()

	caloriesDeps, err := calorieswiring.BuildCalories(cfg)
	if err != nil {
		logger.Error("calories db error", "error", err)
		os.Exit(1)
	}
	defer func() {
		if err := caloriesDeps.Close(); err != nil {
			logger.Error("calories db close error", "error", err)
		}
	}()

	usersDB, err := userswiring.OpenUsersDB(cfg)
	if err != nil {
		logger.Error("users/auth db error", "error", err)
		os.Exit(1)
	}
	defer func() {
		if err := usersDB.Close(); err != nil {
			logger.Error("users/auth db close error", "error", err)
		}
	}()
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

	router := httptransport.NewRouter(handlers, cfg.CORSAllowedOrigin)

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	shutdownDone := make(chan os.Signal, 1)
	signal.Notify(shutdownDone, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-shutdownDone
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		logger.Info("shutdown started")
		if err := srv.Shutdown(ctx); err != nil {
			logger.Error("shutdown error", "error", err)
		}
	}()

	logger.Info("server started", "port", cfg.Port, "env", cfg.Env)
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error("server error", "error", err)
	}
}

func logTursoSelection(logger *slog.Logger, domain string, cfg commonwiring.TursoConfig) {
	selection, err := commonwiring.ResolveTursoSelection(cfg)
	if err != nil {
		logger.Error("turso config error", "domain", domain, "error", err)
		return
	}

	logger.Info("turso connection selected", "domain", domain, "mode", selection.Mode, "scheme", selection.Scheme, "host", selection.Host)
}
