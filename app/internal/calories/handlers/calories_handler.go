package handlers

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/23St/trackstack/internal/calories"
	"github.com/23St/trackstack/internal/calories/components"
	"github.com/23St/trackstack/internal/common/server"
)

// CaloriesHandler handles HTTP requests for the calories dashboard view
type CaloriesHandler struct {
	service calories.Service
}

// NewCaloriesHandler creates a new dashboard handler
func NewCaloriesHandler(service calories.Service) *CaloriesHandler {
	return &CaloriesHandler{service: service}
}

// RenderDashboard renders the main dashboard view
func (h *CaloriesHandler) RenderView(w http.ResponseWriter, r *http.Request) {
	// Get user from context (set by session middleware)
	user, ok := server.GetUserFromContext(r.Context())
	if !ok {
		slog.Error("user not found in context")
		http.Error(w, "User not found in context", http.StatusInternalServerError)
		return
	}

	// Get dashboard summary
	summary, err := h.service.CalculateDailySummary(r.Context(), user.ID, time.Now())
	if err != nil {
		slog.Error("failed to get dashboard summary", "error", err, "user_id", user.ID)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Render HTML template
	if err := components.Dashboard(summary).Render(r.Context(), w); err != nil {
		slog.Error("failed to render dashboard", "error", err)
		http.Error(w, "Failed to render page", http.StatusInternalServerError)
		return
	}
}
