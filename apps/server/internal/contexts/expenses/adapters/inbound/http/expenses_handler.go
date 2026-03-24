package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/Hughzu/trackstack/apps/server/internal/contexts/expenses/application/ports"
	"github.com/Hughzu/trackstack/apps/server/internal/contexts/expenses/domain"
	"github.com/Hughzu/trackstack/apps/server/internal/platform/authcontext"
)

type errorResponse struct {
	Error string `json:"error"`
}

type updateSettingsPayload struct {
	Income      *float64 `json:"income"`
	RatioFund   *int     `json:"ratioFund"`
	RatioFun    *int     `json:"ratioFun"`
	RatioFuture *int     `json:"ratioFuture"`
}

type entryPayload struct {
	Title    string   `json:"title"`
	Amount   *float64 `json:"amount"`
	Category *string  `json:"category"`
	Date     *string  `json:"date"`
}

type templatePayload struct {
	ID       *string  `json:"id"`
	Title    string   `json:"title"`
	Amount   *float64 `json:"amount"`
	Category *string  `json:"category"`
}

type completeChecklistPayload struct {
	ID   string  `json:"id"`
	Date *string `json:"date"`
}

type ExpensesHandler struct {
	settingsUseCase  ports.SettingsUseCase
	entryUseCase     ports.EntryUseCase
	templateUseCase  ports.TemplateUseCase
	sheetUseCase     ports.SheetUseCase
	dashboardUseCase ports.DashboardUseCase
}

func NewExpensesHandler(
	settingsUseCase ports.SettingsUseCase,
	entryUseCase ports.EntryUseCase,
	templateUseCase ports.TemplateUseCase,
	sheetUseCase ports.SheetUseCase,
	dashboardUseCase ports.DashboardUseCase,
) *ExpensesHandler {
	return &ExpensesHandler{
		settingsUseCase:  settingsUseCase,
		entryUseCase:     entryUseCase,
		templateUseCase:  templateUseCase,
		sheetUseCase:     sheetUseCase,
		dashboardUseCase: dashboardUseCase,
	}
}

func (h *ExpensesHandler) GetSettings(w http.ResponseWriter, r *http.Request) {
	userID, ok := h.userID(w, r)
	if !ok {
		return
	}

	view, err := h.settingsUseCase.GetSettings(r.Context(), ports.GetSettingsQuery{UserID: userID})
	if err != nil {
		h.writeError(w, err)
		return
	}

	h.writeJSON(w, http.StatusOK, view)
}

func (h *ExpensesHandler) UpdateSettings(w http.ResponseWriter, r *http.Request) {
	userID, ok := h.userID(w, r)
	if !ok {
		return
	}

	var payload updateSettingsPayload
	if err := decodeJSON(r, &payload); err != nil {
		h.writeJSON(w, http.StatusBadRequest, errorResponse{Error: "Invalid JSON body"})
		return
	}
	if payload.Income == nil || payload.RatioFund == nil || payload.RatioFun == nil || payload.RatioFuture == nil {
		h.writeJSON(w, http.StatusBadRequest, errorResponse{Error: "Missing required fields"})
		return
	}

	settings, err := h.settingsUseCase.UpdateSettings(r.Context(), ports.UpdateSettingsCommand{
		UserID:      userID,
		Income:      *payload.Income,
		RatioFund:   *payload.RatioFund,
		RatioFun:    *payload.RatioFun,
		RatioFuture: *payload.RatioFuture,
	})
	if err != nil {
		h.writeError(w, err)
		return
	}

	h.writeJSON(w, http.StatusOK, settings)
}

func (h *ExpensesHandler) GetDashboard(w http.ResponseWriter, r *http.Request) {
	userID, ok := h.userID(w, r)
	if !ok {
		return
	}

	historyLimit := 50
	if limitValue := r.URL.Query().Get("historyLimit"); limitValue != "" {
		if parsed, err := strconv.Atoi(limitValue); err == nil && parsed > 0 {
			historyLimit = parsed
		}
	}

	dashboard, err := h.dashboardUseCase.GetDashboard(r.Context(), ports.GetDashboardQuery{
		UserID:       userID,
		HistoryLimit: historyLimit,
	})
	if err != nil {
		h.writeError(w, err)
		return
	}

	h.writeJSON(w, http.StatusOK, dashboard)
}

func (h *ExpensesHandler) AddEntry(w http.ResponseWriter, r *http.Request) {
	userID, ok := h.userID(w, r)
	if !ok {
		return
	}

	var payload entryPayload
	if err := decodeJSON(r, &payload); err != nil {
		h.writeJSON(w, http.StatusBadRequest, errorResponse{Error: "Invalid JSON body"})
		return
	}
	if payload.Amount == nil {
		h.writeJSON(w, http.StatusBadRequest, errorResponse{Error: "Missing required fields"})
		return
	}

	entry, err := h.entryUseCase.AddEntry(r.Context(), ports.AddEntryCommand{
		UserID:   userID,
		Title:    payload.Title,
		Amount:   *payload.Amount,
		Category: payload.Category,
		Date:     payload.Date,
	})
	if err != nil {
		h.writeError(w, err)
		return
	}

	h.writeJSON(w, http.StatusCreated, entry)
}

func (h *ExpensesHandler) DeleteEntry(w http.ResponseWriter, r *http.Request, id string) {
	userID, ok := h.userID(w, r)
	if !ok {
		return
	}

	if strings.TrimSpace(id) == "" {
		h.writeJSON(w, http.StatusBadRequest, errorResponse{Error: "Missing expense id"})
		return
	}

	deleted, err := h.entryUseCase.DeleteEntry(r.Context(), ports.DeleteEntryCommand{
		UserID: userID,
		ID:     id,
	})
	if err != nil {
		h.writeError(w, err)
		return
	}
	if !deleted {
		h.writeJSON(w, http.StatusNotFound, errorResponse{Error: "Expense not found"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ExpensesHandler) UpsertChecklist(w http.ResponseWriter, r *http.Request) {
	userID, ok := h.userID(w, r)
	if !ok {
		return
	}

	var payload templatePayload
	if err := decodeJSON(r, &payload); err != nil {
		h.writeJSON(w, http.StatusBadRequest, errorResponse{Error: "Invalid JSON body"})
		return
	}
	if strings.TrimSpace(payload.Title) == "" || payload.Amount == nil {
		h.writeJSON(w, http.StatusBadRequest, errorResponse{Error: "Missing required fields"})
		return
	}

	template, err := h.templateUseCase.UpsertChecklist(r.Context(), ports.UpsertTemplateCommand{
		ID:       payload.ID,
		UserID:   userID,
		Title:    payload.Title,
		Amount:   *payload.Amount,
		Category: payload.Category,
	})
	if err != nil {
		h.writeError(w, err)
		return
	}

	h.writeJSON(w, http.StatusOK, template)
}

func (h *ExpensesHandler) DeleteChecklist(w http.ResponseWriter, r *http.Request, id string) {
	userID, ok := h.userID(w, r)
	if !ok {
		return
	}

	if strings.TrimSpace(id) == "" {
		h.writeJSON(w, http.StatusBadRequest, errorResponse{Error: "Missing template id"})
		return
	}

	deleted, err := h.templateUseCase.DeleteChecklist(r.Context(), ports.DeleteTemplateCommand{
		UserID: userID,
		ID:     id,
	})
	if err != nil {
		h.writeError(w, err)
		return
	}
	if !deleted {
		h.writeJSON(w, http.StatusNotFound, errorResponse{Error: "Template not found"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ExpensesHandler) CompleteChecklistItem(w http.ResponseWriter, r *http.Request) {
	userID, ok := h.userID(w, r)
	if !ok {
		return
	}

	var payload completeChecklistPayload
	if err := decodeJSON(r, &payload); err != nil {
		h.writeJSON(w, http.StatusBadRequest, errorResponse{Error: "Invalid JSON body"})
		return
	}
	if strings.TrimSpace(payload.ID) == "" {
		h.writeJSON(w, http.StatusBadRequest, errorResponse{Error: "Missing checklist item id"})
		return
	}

	entry, err := h.templateUseCase.CompleteChecklistItem(r.Context(), ports.CompleteChecklistItemCommand{
		UserID: userID,
		ID:     payload.ID,
		Date:   payload.Date,
	})
	if err != nil {
		h.writeError(w, err)
		return
	}

	h.writeJSON(w, http.StatusCreated, entry)
}

func (h *ExpensesHandler) UpsertRecurring(w http.ResponseWriter, r *http.Request) {
	userID, ok := h.userID(w, r)
	if !ok {
		return
	}

	var payload templatePayload
	if err := decodeJSON(r, &payload); err != nil {
		h.writeJSON(w, http.StatusBadRequest, errorResponse{Error: "Invalid JSON body"})
		return
	}
	if strings.TrimSpace(payload.Title) == "" || payload.Amount == nil {
		h.writeJSON(w, http.StatusBadRequest, errorResponse{Error: "Missing required fields"})
		return
	}

	template, err := h.templateUseCase.UpsertRecurring(r.Context(), ports.UpsertTemplateCommand{
		ID:       payload.ID,
		UserID:   userID,
		Title:    payload.Title,
		Amount:   *payload.Amount,
		Category: payload.Category,
	})
	if err != nil {
		h.writeError(w, err)
		return
	}

	h.writeJSON(w, http.StatusOK, template)
}

func (h *ExpensesHandler) DeleteRecurring(w http.ResponseWriter, r *http.Request, id string) {
	userID, ok := h.userID(w, r)
	if !ok {
		return
	}

	if strings.TrimSpace(id) == "" {
		h.writeJSON(w, http.StatusBadRequest, errorResponse{Error: "Missing template id"})
		return
	}

	deleted, err := h.templateUseCase.DeleteRecurring(r.Context(), ports.DeleteTemplateCommand{
		UserID: userID,
		ID:     id,
	})
	if err != nil {
		h.writeError(w, err)
		return
	}
	if !deleted {
		h.writeJSON(w, http.StatusNotFound, errorResponse{Error: "Template not found"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ExpensesHandler) CloseSheet(w http.ResponseWriter, r *http.Request) {
	userID, ok := h.userID(w, r)
	if !ok {
		return
	}

	sheet, err := h.sheetUseCase.CloseSheet(r.Context(), ports.CloseSheetCommand{UserID: userID})
	if err != nil {
		h.writeError(w, err)
		return
	}

	h.writeJSON(w, http.StatusOK, sheet)
}

func (h *ExpensesHandler) userID(w http.ResponseWriter, r *http.Request) (string, bool) {
	userID, ok := authcontext.GetUserID(r.Context())
	if !ok || userID == "" {
		h.writeJSON(w, http.StatusUnauthorized, errorResponse{Error: "Unauthorized"})
		return "", false
	}

	return userID, true
}

func decodeJSON(r *http.Request, dst any) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(dst)
}

func (h *ExpensesHandler) writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func (h *ExpensesHandler) writeError(w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	if errors.Is(err, domain.ErrInvalidInput) {
		status = http.StatusBadRequest
	}
	if errors.Is(err, domain.ErrNotFound) {
		status = http.StatusNotFound
	}

	h.writeJSON(w, status, errorResponse{Error: err.Error()})
}
