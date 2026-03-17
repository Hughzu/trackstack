package http

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	heatservices "github.com/Hughzu/trackstack/apps/server/internal/contexts/heat/application/services"
	heatdomain "github.com/Hughzu/trackstack/apps/server/internal/contexts/heat/domain"
	"github.com/go-chi/chi/v5"
)

type ListRefillsFunc func(context.Context, heatservices.ListRefillsRequest) ([]heatdomain.Refill, error)
type GetDashboardFunc func(context.Context, heatservices.GetDashboardRequest) (heatservices.DashboardViewModel, error)
type CreateRefillFunc func(context.Context, heatservices.CreateRefillRequest) (heatdomain.Refill, error)
type DeleteRefillFunc func(context.Context, heatservices.DeleteRefillRequest) (bool, error)

type Handler struct {
	listRefills       ListRefillsFunc
	getDashboard      GetDashboardFunc
	createRefill      CreateRefillFunc
	deleteRefill      DeleteRefillFunc
	requireAuthUserID func(http.ResponseWriter, *http.Request) (string, bool)
	writeJSON         func(http.ResponseWriter, int, any)
	decodeJSON        func(*http.Request, any) error
}

type Deps struct {
	ListRefills       ListRefillsFunc
	GetDashboard      GetDashboardFunc
	CreateRefill      CreateRefillFunc
	DeleteRefill      DeleteRefillFunc
	RequireAuthUserID func(http.ResponseWriter, *http.Request) (string, bool)
	WriteJSON         func(http.ResponseWriter, int, any)
	DecodeJSON        func(*http.Request, any) error
}

type createRefillPayload struct {
	Date        string   `json:"date"`
	WeightKg    float64  `json:"weightKg"`
	Bags        int      `json:"bags"`
	Temperature *float64 `json:"temperature"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func NewHandler(deps Deps) *Handler {
	return &Handler{
		listRefills:       deps.ListRefills,
		getDashboard:      deps.GetDashboard,
		createRefill:      deps.CreateRefill,
		deleteRefill:      deps.DeleteRefill,
		requireAuthUserID: deps.RequireAuthUserID,
		writeJSON:         deps.WriteJSON,
		decodeJSON:        deps.DecodeJSON,
	}
}

func MountRoutes(r chi.Router, handler *Handler) {
	r.Get("/dashboard", handler.GetDashboard)
	r.Get("/refills", handler.ListRefills)
	r.Post("/refills", handler.CreateRefill)
	r.Delete("/refills", handler.DeleteRefill)
	r.Delete("/refills/{id}", handler.DeleteRefill)
}

func (h *Handler) ListRefills(w http.ResponseWriter, r *http.Request) {
	userID, ok := h.requireAuthUserID(w, r)
	if !ok {
		return
	}

	refills, err := h.listRefills(r.Context(), heatservices.ListRefillsRequest{
		UserID: userID,
		From:   r.URL.Query().Get("from"),
		To:     r.URL.Query().Get("to"),
	})
	if err != nil {
		h.writeError(w, err)
		return
	}

	h.writeJSON(w, http.StatusOK, refills)
}

func (h *Handler) GetDashboard(w http.ResponseWriter, r *http.Request) {
	userID, ok := h.requireAuthUserID(w, r)
	if !ok {
		return
	}

	page := 1
	if pageValue := r.URL.Query().Get("page"); pageValue != "" {
		if parsed, err := strconv.Atoi(pageValue); err == nil && parsed > 0 {
			page = parsed
		}
	}

	limit := 20
	if limitValue := r.URL.Query().Get("limit"); limitValue != "" {
		if parsed, err := strconv.Atoi(limitValue); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	dashboard, err := h.getDashboard(r.Context(), heatservices.GetDashboardRequest{
		UserID: userID,
		Page:   page,
		Limit:  limit,
	})
	if err != nil {
		h.writeError(w, err)
		return
	}

	h.writeJSON(w, http.StatusOK, dashboard)
}

func (h *Handler) CreateRefill(w http.ResponseWriter, r *http.Request) {
	userID, ok := h.requireAuthUserID(w, r)
	if !ok {
		return
	}

	var payload createRefillPayload
	if err := h.decodeJSON(r, &payload); err != nil {
		h.writeJSON(w, http.StatusBadRequest, errorResponse{Error: "Invalid JSON body"})
		return
	}

	refill, err := h.createRefill(r.Context(), heatservices.CreateRefillRequest{
		UserID:      userID,
		Date:        payload.Date,
		WeightKg:    payload.WeightKg,
		Bags:        payload.Bags,
		Temperature: payload.Temperature,
	})
	if err != nil {
		h.writeError(w, err)
		return
	}

	h.writeJSON(w, http.StatusCreated, refill)
}

func (h *Handler) DeleteRefill(w http.ResponseWriter, r *http.Request) {
	userID, ok := h.requireAuthUserID(w, r)
	if !ok {
		return
	}

	id := chi.URLParam(r, "id")
	if id == "" {
		id = r.URL.Query().Get("id")
	}
	if id == "" {
		h.writeJSON(w, http.StatusBadRequest, errorResponse{Error: "Missing refill id"})
		return
	}

	deleted, err := h.deleteRefill(r.Context(), heatservices.DeleteRefillRequest{
		UserID: userID,
		ID:     id,
	})
	if err != nil {
		h.writeError(w, err)
		return
	}
	if !deleted {
		h.writeJSON(w, http.StatusNotFound, errorResponse{Error: "Refill not found"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) writeError(w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	if errors.Is(err, heatdomain.ErrInvalidInput) {
		status = http.StatusBadRequest
	}

	h.writeJSON(w, status, errorResponse{Error: err.Error()})
}
