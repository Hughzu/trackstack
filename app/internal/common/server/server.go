package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/23St/trackstack/internal/common/db"
)

// ViewHandler defines the interface for view rendering
type ViewHandler interface {
	RenderView(w http.ResponseWriter, r *http.Request)
}

// Handlers groups all HTTP handlers for dependency injection.
// TODO meh, when split to microservices it will break
type Handlers struct {
	Calories ViewHandler
	//Expenses ViewHandler
	//Metrics ViewHandler
}

// Server represents the HTTP server
type Server struct {
	httpServer *http.Server
	mux        *http.ServeMux
	db         *db.DB
	handlers   Handlers
}

// NewServer creates a new HTTP server
func NewServer(port string, database *db.DB, handlers Handlers) *Server {
	mux := http.NewServeMux()

	s := &Server{
		httpServer: &http.Server{
			Addr:         fmt.Sprintf(":%s", port),
			Handler:      SessionMiddleware(database)(mux),
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
		mux:      mux,
		db:       database,
		handlers: handlers,
	}

	// Register routes
	s.registerRoutes()

	return s
}

// registerRoutes sets up all HTTP routes
func (s *Server) registerRoutes() {
	// Calories (root route)
	// TODO routes inside modules ?
	s.mux.HandleFunc("/", s.handlers.Calories.RenderView)

	// API routes
	s.mux.HandleFunc("/health", s.handleHealth)
	s.mux.HandleFunc("/api/session", s.handleSessionInfo)

	// Static files
	s.mux.Handle("/static/", http.StripPrefix("/static/",
		http.FileServer(http.Dir("./static"))))
}

// handleHealth returns a simple health check response
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	response := map[string]string{
		"status": "ok",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// handleSessionInfo returns current user session information
func (s *Server) handleSessionInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user, ok := GetUserFromContext(r.Context())
	if !ok {
		http.Error(w, "No user in context", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// Start begins listening for HTTP requests
func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
