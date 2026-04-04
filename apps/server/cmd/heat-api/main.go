package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Hughzu/trackstack/apps/server/internal/app/heatapi"
	"github.com/joho/godotenv"
)

func init() {
	_ = godotenv.Load()
}

func main() {
	runtime, err := heatapi.NewRuntime()
	if err != nil {
		_, _ = os.Stderr.WriteString("Startup error : " + err.Error() + "\n")
		os.Exit(1)
	}

	defer func() {
		_ = runtime.Close()
	}()

	srv := &http.Server{
		Addr:         ":" + runtime.Config.Port,
		Handler:      runtime.Handler,
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

		runtime.Logger.Info("Heat service shutdown started")
		if err := srv.Shutdown(ctx); err != nil {
			runtime.Logger.Error("Heat service shutdown error", "error", err)
		}
	}()

	runtime.Logger.Info("Heat service started", "port", runtime.Config.Port, "env", runtime.Config.Env)

	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		runtime.Logger.Error("Heat service error", "error", err)
		os.Exit(1)
	}
}
