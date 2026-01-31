package handlers

import (
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/23St/trackstack/internal/calories"
	"github.com/23St/trackstack/internal/calories/components"
	"github.com/23St/trackstack/internal/common/server"
)

// Handler handles HTTP requests for the calories module
type Handler struct {
	service calories.Service
}

// NewHandler creates a new handler
func NewHandler(service calories.Service) *Handler {
	return &Handler{service: service}
}

// RegisterRoutes registers all routes for the calories module
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/calories/", h.Dashboard)
	mux.HandleFunc("/calories/log", h.LogMeal)
}

// Dashboard renders the main calories dashboard view
func (h *Handler) Dashboard(w http.ResponseWriter, r *http.Request) {
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

// LogMeal handles POST /calories/log for manual meal entry
func (h *Handler) LogMeal(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get user from context
	user, ok := server.GetUserFromContext(r.Context())
	if !ok {
		slog.Error("user not found in context")
		http.Error(w, "User not found in context", http.StatusInternalServerError)
		return
	}

	// Parse form values
	name := strings.TrimSpace(r.FormValue("name"))
	kcals, _ := strconv.Atoi(r.FormValue("calories"))
	protein, _ := strconv.Atoi(r.FormValue("protein"))
	carbs, _ := strconv.Atoi(r.FormValue("carbs"))
	fat, _ := strconv.Atoi(r.FormValue("fat"))

	// Validate required fields
	if name == "" || kcals <= 0 {
		http.Error(w, "Meal name and calories are required", http.StatusBadRequest)
		return
	}

	// Log the meal
	summary, err := h.service.LogMeal(r.Context(), user.ID, name, kcals, protein, carbs, fat)
	if err != nil {
		slog.Error("failed to log meal", "error", err, "user_id", user.ID)
		http.Error(w, "Failed to log meal", http.StatusInternalServerError)
		return
	}

	// Render only the dashboard metrics for HTMX swap
	if err := components.DashboardMetrics(summary).Render(r.Context(), w); err != nil {
		slog.Error("failed to render dashboard metrics", "error", err)
		http.Error(w, "Failed to render", http.StatusInternalServerError)
		return
	}

	// Render the recent meals list as an OOB swap
	if err := components.RecentMealsList(summary.RecentMeals, true).Render(r.Context(), w); err != nil {
		slog.Error("failed to render recent meals list", "error", err)
		// Don't fail the request if just the OOB update fails
	}
}
