package httptransport_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	heatapp "github.com/Hughzu/trackstack/apps/server/internal/contexts/heat/application"
	"github.com/Hughzu/trackstack/apps/server/internal/modules/auth"
	"github.com/Hughzu/trackstack/apps/server/internal/modules/calories"
	"github.com/Hughzu/trackstack/apps/server/internal/modules/expenses"
	"github.com/Hughzu/trackstack/apps/server/internal/modules/users"
	httptransport "github.com/Hughzu/trackstack/apps/server/internal/transport/http"
	"github.com/go-chi/chi/v5"
)

type mockExpensesStore struct {
	settingsUpdated  bool
	checklistCreated bool
	recurringCreated bool
	checklistDeleted string
	recurringDeleted string
	settings         expenses.Settings
	checklist        []expenses.Template
	recurring        []expenses.Template
}

func newMockExpensesStore() *mockExpensesStore {
	return &mockExpensesStore{
		settings: expenses.Settings{
			ID:          "settings-1",
			UserID:      "test-user",
			Income:      2000,
			RatioFund:   60,
			RatioFun:    20,
			RatioFuture: 20,
			CreatedAt:   "2026-03-10T00:00:00Z",
			UpdatedAt:   "2026-03-10T00:00:00Z",
		},
	}
}

func (m *mockExpensesStore) GetSettings(ctx context.Context, userID string) (expenses.Settings, error) {
	return m.settings, nil
}

func (m *mockExpensesStore) CreateSettings(ctx context.Context, settings expenses.Settings) error {
	m.settings = settings
	return nil
}

func (m *mockExpensesStore) UpdateSettings(ctx context.Context, settings expenses.Settings) error {
	m.settingsUpdated = true
	m.settings = settings
	return nil
}

func (m *mockExpensesStore) GetChecklistTemplates(ctx context.Context, userID string) ([]expenses.Template, error) {
	return m.checklist, nil
}

func (m *mockExpensesStore) GetChecklistTemplate(ctx context.Context, templateID string, userID string) (expenses.Template, error) {
	return expenses.Template{}, expenses.ErrNotFound
}

func (m *mockExpensesStore) CreateChecklistTemplate(ctx context.Context, template expenses.Template) error {
	m.checklistCreated = true
	m.checklist = append(m.checklist, template)
	return nil
}

func (m *mockExpensesStore) UpdateChecklistTemplate(ctx context.Context, template expenses.Template) error {
	return nil
}

func (m *mockExpensesStore) DeleteChecklistTemplate(ctx context.Context, templateID string, userID string) (bool, error) {
	m.checklistDeleted = templateID
	return true, nil
}

func (m *mockExpensesStore) GetRecurringTemplates(ctx context.Context, userID string) ([]expenses.Template, error) {
	return m.recurring, nil
}

func (m *mockExpensesStore) GetRecurringTemplate(ctx context.Context, templateID string, userID string) (expenses.Template, error) {
	return expenses.Template{}, expenses.ErrNotFound
}

func (m *mockExpensesStore) CreateRecurringTemplate(ctx context.Context, template expenses.Template) error {
	m.recurringCreated = true
	m.recurring = append(m.recurring, template)
	return nil
}

func (m *mockExpensesStore) UpdateRecurringTemplate(ctx context.Context, template expenses.Template) error {
	return nil
}

func (m *mockExpensesStore) DeleteRecurringTemplate(ctx context.Context, templateID string, userID string) (bool, error) {
	m.recurringDeleted = templateID
	return true, nil
}

func (m *mockExpensesStore) GetOpenSheet(ctx context.Context, userID string) (*expenses.Sheet, error) {
	return &expenses.Sheet{ID: "sheet-1", UserID: userID, PeriodKey: "2026-03"}, nil
}

func (m *mockExpensesStore) GetLatestSheet(ctx context.Context, userID string) (*expenses.Sheet, error) {
	return &expenses.Sheet{ID: "sheet-1", UserID: userID, PeriodKey: "2026-03"}, nil
}

func (m *mockExpensesStore) CreateSheet(ctx context.Context, sheet expenses.Sheet) error {
	return nil
}

func (m *mockExpensesStore) UpdateSheet(ctx context.Context, sheet expenses.Sheet) error {
	return nil
}

func (m *mockExpensesStore) GetChecklistItem(ctx context.Context, itemID string, userID string) (expenses.ChecklistItem, error) {
	return expenses.ChecklistItem{}, expenses.ErrNotFound
}

func (m *mockExpensesStore) GetPendingChecklistItems(ctx context.Context, sheetID string) ([]expenses.ChecklistItem, error) {
	return nil, nil
}

func (m *mockExpensesStore) CreateChecklistItem(ctx context.Context, item expenses.ChecklistItem) error {
	return nil
}

func (m *mockExpensesStore) UpdateChecklistItem(ctx context.Context, item expenses.ChecklistItem) error {
	return nil
}

func (m *mockExpensesStore) UpdateChecklistItemsByTemplate(ctx context.Context, templateID string, title string, amount float64, category expenses.Category) error {
	return nil
}

func (m *mockExpensesStore) DeletePendingChecklistItemsByTemplate(ctx context.Context, templateID string, userID string) error {
	return nil
}

func (m *mockExpensesStore) CreateExpense(ctx context.Context, entry expenses.Entry) error {
	return nil
}

func (m *mockExpensesStore) DeleteExpense(ctx context.Context, entryID string, userID string) (bool, error) {
	return true, nil
}

func (m *mockExpensesStore) GetSheetHistory(ctx context.Context, sheetID string) ([]expenses.Entry, error) {
	return nil, nil
}

func (m *mockExpensesStore) GetRecentSheetHistory(ctx context.Context, sheetID string, limit int, offset int) ([]expenses.Entry, error) {
	return []expenses.Entry{}, nil
}

func (m *mockExpensesStore) GetTotalSpentBySheet(ctx context.Context, sheetID string) (float64, error) {
	return 0, nil
}

func (m *mockExpensesStore) GetSpentByCategory(ctx context.Context, sheetID string) (map[expenses.Category]float64, error) {
	return map[expenses.Category]float64{}, nil
}

func setupExpensesTestRouter(store *mockExpensesStore) *chi.Mux {
	authService := auth.NewService(&testAuthStore{}, testAuthConfig())
	usersService := users.NewService(&testUsersStore{})
	expensesService := expenses.NewService(store)

	handlers := httptransport.NewHandlers(httptransport.Deps{
		AuthService:        authService,
		UsersService:       usersService,
		CaloriesService:    calories.NewService(nil),
		HeatService:        &heatapp.Service{},
		ExpensesService:    expensesService,
		AuthCookieName:     testCookieName,
		AuthCookieSecure:   false,
		AuthCookieSameSite: "lax",
	})

	return httptransport.NewRouter(handlers, "*").(*chi.Mux)
}

func TestExpensesUpdateSettingsAPI_JSON(t *testing.T) {
	store := newMockExpensesStore()
	router := setupExpensesTestRouter(store)

	body, err := json.Marshal(map[string]any{
		"income":      2500,
		"ratioFund":   50,
		"ratioFun":    25,
		"ratioFuture": 25,
	})
	if err != nil {
		t.Fatalf("marshal payload: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/expenses/settings", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{Name: testCookieName, Value: testSessionToken})

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: body: %s", rr.Code, rr.Body.String())
	}
	if !store.settingsUpdated {
		t.Fatalf("expected settings to be updated")
	}
}

func TestExpensesUpsertChecklistAPI_JSON(t *testing.T) {
	store := newMockExpensesStore()
	router := setupExpensesTestRouter(store)

	body, err := json.Marshal(map[string]any{
		"title":    "Netflix",
		"amount":   19.99,
		"category": "fun",
	})
	if err != nil {
		t.Fatalf("marshal payload: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/expenses/checklists", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{Name: testCookieName, Value: testSessionToken})

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: body: %s", rr.Code, rr.Body.String())
	}
	if !store.checklistCreated {
		t.Fatalf("expected checklist template to be created")
	}
}

func TestExpensesUpsertRecurringAPI_JSON(t *testing.T) {
	store := newMockExpensesStore()
	router := setupExpensesTestRouter(store)

	body, err := json.Marshal(map[string]any{
		"title":    "Rent",
		"amount":   900,
		"category": "fund",
	})
	if err != nil {
		t.Fatalf("marshal payload: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/expenses/recurring", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{Name: testCookieName, Value: testSessionToken})

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: body: %s", rr.Code, rr.Body.String())
	}
	if !store.recurringCreated {
		t.Fatalf("expected recurring template to be created")
	}
}

func TestExpensesAddExpenseAPI_FormRejected(t *testing.T) {
	store := newMockExpensesStore()
	router := setupExpensesTestRouter(store)

	form := url.Values{}
	form.Set("title", "Groceries")
	form.Set("amount", "42.50")
	form.Set("category", "fund")

	req := httptest.NewRequest(http.MethodPost, "/api/expenses/entries", bytes.NewBufferString(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(&http.Cookie{Name: testCookieName, Value: testSessionToken})

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400 for non-JSON request, got %d: body: %s", rr.Code, rr.Body.String())
	}
}

func TestExpensesLegacyAliasesReturnNotFound(t *testing.T) {
	store := newMockExpensesStore()
	router := setupExpensesTestRouter(store)

	tests := []struct {
		name   string
		method string
		path   string
	}{
		{name: "expense post alias removed", method: http.MethodPost, path: "/api/expenses/expense"},
		{name: "expense delete alias removed", method: http.MethodDelete, path: "/api/expenses/expense?id=entry-1"},
		{name: "checklist post alias removed", method: http.MethodPost, path: "/api/expenses/checklist"},
		{name: "checklist delete alias removed", method: http.MethodDelete, path: "/api/expenses/checklist?id=template-1"},
		{name: "checklist complete alias removed", method: http.MethodPost, path: "/api/expenses/checklist/complete"},
		{name: "close alias removed", method: http.MethodPost, path: "/api/expenses/close"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			req.AddCookie(&http.Cookie{Name: testCookieName, Value: testSessionToken})

			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			if rr.Code != http.StatusNotFound {
				t.Fatalf("expected status 404, got %d: body: %s", rr.Code, rr.Body.String())
			}
		})
	}
}
