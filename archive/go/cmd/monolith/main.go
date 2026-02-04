package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/23St/trackstack/internal/calories"
	"github.com/23St/trackstack/internal/calories/handlers"
	"github.com/23St/trackstack/internal/common/db"
	"github.com/23St/trackstack/internal/common/server"
)

func main() {
	// Initialize structured logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(logger)

	// Initialize database
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./data/trackstack.db"
	}

	database, err := db.New(dbPath)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer func(database *db.DB) {
		err := database.Close()
		if err != nil {
			slog.Error("failed to close database connection", "error", err)
			os.Exit(1)
		}
	}(database)

	if err := database.InitSchema(); err != nil {
		slog.Error("failed to initialize schema", "error", err)
		os.Exit(1)
	}

	slog.Info("database initialized", "path", dbPath)

	// Initialize modules
	caloriesMod := calories.NewModule(database)

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Create server
	srv := server.NewServer(port, database)

	// Create handler and register routes
	caloriesHandler := handlers.NewHandler(caloriesMod.Service)
	srv.RegisterRoutes(caloriesHandler.RegisterRoutes)

	// Redirect root to calories module
	srv.RegisterRoutes(func(mux *http.ServeMux) {
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/" {
				http.Redirect(w, r, "/calories/", http.StatusMovedPermanently)
				return
			}
			http.NotFound(w, r)
		})
	})

	// Start server in goroutine
	go func() {
		slog.Info("starting server", "port", port)
		if err := srv.Start(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("shutting down server")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("server forced to shutdown", "error", err)
		os.Exit(1)
	}

	slog.Info("server exited")
}
