package httptransport

import (
	"errors"
	"net/http"

	"github.com/Hughzu/trackstack/apps/server/internal/modules/expenses"
)

type ExpensesHandler struct {
	svc *expenses.Service
}

func (h *ExpensesHandler) GetSettings(w http.ResponseWriter, r *http.Request) {
	userID, ok := requireAuthUserID(w, r)
	if !ok {
		return
	}
	settings, err := h.svc.GetSettings(r.Context(), expenses.GetSettingsRequest{
		UserID: userID,
	})
	if err != nil {
		writeExpensesError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, settings)
}

func (h *ExpensesHandler) UpdateSettings(w http.ResponseWriter, r *http.Request) {
	userID, ok := requireAuthUserID(w, r)
	if !ok {
		return
	}
	var payload struct {
		Income      *float64 `json:"income"`
		RatioFund   *int     `json:"ratioFund"`
		RatioFun    *int     `json:"ratioFun"`
		RatioFuture *int     `json:"ratioFuture"`
	}
	if err := decodeJSON(r, &payload); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "Invalid JSON body"})
		return
	}

	settings, err := h.svc.UpdateSettings(r.Context(), expenses.UpdateSettingsRequest{
		UserID:      userID,
		Income:      payload.Income,
		RatioFund:   payload.RatioFund,
		RatioFun:    payload.RatioFun,
		RatioFuture: payload.RatioFuture,
	})
	if err != nil {
		writeExpensesError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, settings)
}

func (h *ExpensesHandler) GetCurrentSheet(w http.ResponseWriter, r *http.Request) {
	userID, ok := requireAuthUserID(w, r)
	if !ok {
		return
	}
	dashboard, err := h.svc.GetDashboard(r.Context(), expenses.GetCurrentSheetRequest{
		UserID: userID,
	})
	if err != nil {
		writeExpensesError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, dashboard)
}

func (h *ExpensesHandler) AddExpense(w http.ResponseWriter, r *http.Request) {
	userID, ok := requireAuthUserID(w, r)
	if !ok {
		return
	}
	var payload struct {
		Title    string  `json:"title"`
		Amount   float64 `json:"amount"`
		Category *string `json:"category"`
		Date     *string `json:"date"`
	}
	if err := decodeJSON(r, &payload); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "Invalid JSON body"})
		return
	}

	entry, err := h.svc.AddExpense(r.Context(), expenses.AddExpenseRequest{
		UserID:   userID,
		Title:    payload.Title,
		Amount:   payload.Amount,
		Category: payload.Category,
		Date:     payload.Date,
	})
	if err != nil {
		writeExpensesError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, entry)
}

func (h *ExpensesHandler) DeleteExpense(w http.ResponseWriter, r *http.Request) {
	userID, ok := requireAuthUserID(w, r)
	if !ok {
		return
	}
	id := r.URL.Query().Get("id")
	if id == "" {
		id = extractIDFromBody(r)
	}

	if id == "" {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "Missing expense id"})
		return
	}

	deleted, err := h.svc.DeleteExpense(r.Context(), expenses.DeleteExpenseRequest{
		UserID: userID,
		ID:     id,
	})
	if err != nil {
		writeExpensesError(w, err)
		return
	}
	if !deleted {
		writeJSON(w, http.StatusNotFound, errorResponse{Error: "Expense not found"})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *ExpensesHandler) UpsertChecklist(w http.ResponseWriter, r *http.Request) {
	userID, ok := requireAuthUserID(w, r)
	if !ok {
		return
	}
	var payload struct {
		ID       *string `json:"id"`
		Title    string  `json:"title"`
		Amount   float64 `json:"amount"`
		Category *string `json:"category"`
	}
	if err := decodeJSON(r, &payload); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "Invalid JSON body"})
		return
	}

	template, err := h.svc.UpsertChecklistTemplate(r.Context(), expenses.UpsertTemplateRequest{
		ID:       payload.ID,
		UserID:   userID,
		Title:    payload.Title,
		Amount:   payload.Amount,
		Category: payload.Category,
	})
	if err != nil {
		writeExpensesError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, template)
}

func (h *ExpensesHandler) DeleteChecklist(w http.ResponseWriter, r *http.Request) {
	userID, ok := requireAuthUserID(w, r)
	if !ok {
		return
	}
	id := r.URL.Query().Get("id")
	if id == "" {
		id = extractIDFromBody(r)
	}
	if id == "" {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "Missing template id"})
		return
	}

	deleted, err := h.svc.DeleteChecklistTemplate(r.Context(), expenses.DeleteTemplateRequest{
		UserID: userID,
		ID:     id,
	})
	if err != nil {
		writeExpensesError(w, err)
		return
	}
	if !deleted {
		writeJSON(w, http.StatusNotFound, errorResponse{Error: "Template not found"})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *ExpensesHandler) CompleteChecklistItem(w http.ResponseWriter, r *http.Request) {
	userID, ok := requireAuthUserID(w, r)
	if !ok {
		return
	}
	var payload struct {
		ID   string  `json:"id"`
		Date *string `json:"date"`
	}
	if err := decodeJSON(r, &payload); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "Invalid JSON body"})
		return
	}

	entry, err := h.svc.CompleteChecklistItem(r.Context(), expenses.CompleteChecklistItemRequest{
		ID:     payload.ID,
		UserID: userID,
		Date:   payload.Date,
	})
	if err != nil {
		writeExpensesError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, entry)
}

func (h *ExpensesHandler) UpsertRecurring(w http.ResponseWriter, r *http.Request) {
	userID, ok := requireAuthUserID(w, r)
	if !ok {
		return
	}
	var payload struct {
		ID       *string `json:"id"`
		Title    string  `json:"title"`
		Amount   float64 `json:"amount"`
		Category *string `json:"category"`
	}
	if err := decodeJSON(r, &payload); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "Invalid JSON body"})
		return
	}

	template, err := h.svc.UpsertRecurringTemplate(r.Context(), expenses.UpsertTemplateRequest{
		ID:       payload.ID,
		UserID:   userID,
		Title:    payload.Title,
		Amount:   payload.Amount,
		Category: payload.Category,
	})
	if err != nil {
		writeExpensesError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, template)
}

func (h *ExpensesHandler) DeleteRecurring(w http.ResponseWriter, r *http.Request) {
	userID, ok := requireAuthUserID(w, r)
	if !ok {
		return
	}
	id := r.URL.Query().Get("id")
	if id == "" {
		id = extractIDFromBody(r)
	}
	if id == "" {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "Missing template id"})
		return
	}

	deleted, err := h.svc.DeleteRecurringTemplate(r.Context(), expenses.DeleteTemplateRequest{
		UserID: userID,
		ID:     id,
	})
	if err != nil {
		writeExpensesError(w, err)
		return
	}
	if !deleted {
		writeJSON(w, http.StatusNotFound, errorResponse{Error: "Template not found"})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *ExpensesHandler) CloseSheet(w http.ResponseWriter, r *http.Request) {
	userID, ok := requireAuthUserID(w, r)
	if !ok {
		return
	}
	sheet, err := h.svc.CloseSheet(r.Context(), expenses.CloseSheetRequest{
		UserID: userID,
	})
	if err != nil {
		writeExpensesError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, sheet)
}

func writeExpensesError(w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	if errors.Is(err, expenses.ErrInvalidInput) {
		status = http.StatusBadRequest
	}
	writeJSON(w, status, errorResponse{Error: err.Error()})
}
