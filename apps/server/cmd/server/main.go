package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Hughzu/trackstack/apps/server/internal/core/config"
	"github.com/Hughzu/trackstack/apps/server/internal/core/logging"
	httptransport "github.com/Hughzu/trackstack/apps/server/internal/transport/http"
	expenseswiring "github.com/Hughzu/trackstack/apps/server/internal/wiring/expenses"
	heatwiring "github.com/Hughzu/trackstack/apps/server/internal/wiring/heat"

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

	handlers := httptransport.NewHandlers(httptransport.Deps{
		HeatService:     heatDeps.Service,
		ExpensesService: expensesDeps.Service,
		HardcodedUserID: cfg.HardcodedUserID,
	})

	router := httptransport.NewRouter(handlers)

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
