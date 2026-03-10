package httptransport_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/Hughzu/trackstack/apps/server/internal/modules/auth"
	"github.com/Hughzu/trackstack/apps/server/internal/modules/calories"
	"github.com/Hughzu/trackstack/apps/server/internal/modules/expenses"
	"github.com/Hughzu/trackstack/apps/server/internal/modules/heat"
	"github.com/Hughzu/trackstack/apps/server/internal/modules/users"
	httptransport "github.com/Hughzu/trackstack/apps/server/internal/transport/http"
	"github.com/go-chi/chi/v5"
)

// mockCaloriesStore satisfies calories.Store
type mockCaloriesStore struct {
	logAdded bool
}

func (m *mockCaloriesStore) GetTarget(ctx context.Context, userID string) (calories.Target, error) {
	return calories.Target{}, nil
}
func (m *mockCaloriesStore) CreateTarget(ctx context.Context, target calories.Target) error {
	return nil
}
func (m *mockCaloriesStore) UpdateTarget(ctx context.Context, target calories.Target) error {
	return nil
}
func (m *mockCaloriesStore) CreateLog(ctx context.Context, log calories.Log) error {
	m.logAdded = true
	return nil
}
func (m *mockCaloriesStore) DeleteLog(ctx context.Context, userID string, id string) (bool, error) {
	return true, nil
}
func (m *mockCaloriesStore) GetSummaryByRange(ctx context.Context, userID string, from string, to string) (calories.Summary, error) {
	return calories.Summary{}, nil
}
func (m *mockCaloriesStore) GetLogsByRange(ctx context.Context, userID string, from string, to string) ([]calories.Log, error) {
	return nil, nil
}
func (m *mockCaloriesStore) GetRecentLogs(ctx context.Context, userID string, limit int) ([]calories.Log, error) {
	return nil, nil
}

func setupTestRouter(cStore *mockCaloriesStore) *chi.Mux {
	authService := auth.NewService(&testAuthStore{}, testAuthConfig())
	usersService := users.NewService(&testUsersStore{})
	calService := calories.NewService(cStore)

	handlers := httptransport.NewHandlers(httptransport.Deps{
		AuthService:        authService,
		UsersService:       usersService,
		CaloriesService:    calService,
		HeatService:        &heat.Service{},
		ExpensesService:    &expenses.Service{},
		AuthCookieName:     testCookieName,
		AuthCookieSecure:   false,
		AuthCookieSameSite: "lax",
	})

	return httptransport.NewRouter(handlers, "*").(*chi.Mux)
}

func TestCaloriesAddLogAPI_JSON(t *testing.T) {
	store := &mockCaloriesStore{}
	router := setupTestRouter(store)

	payload := map[string]interface{}{
		"calories": 500,
		"protein":  30,
		"title":    "Lunch",
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/api/calories/log", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{Name: testCookieName, Value: testSessionToken})

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d: body: %s", rr.Code, rr.Body.String())
	}
	if !store.logAdded {
		t.Fatalf("expected log to be added to store")
	}
}

func TestCaloriesAddLogAPI_Form(t *testing.T) {
	store := &mockCaloriesStore{}
	router := setupTestRouter(store)

	form := url.Values{}
	form.Add("calories", "600")
	form.Add("protein", "40")
	form.Add("title", "Dinner")

	req := httptest.NewRequest(http.MethodPost, "/api/calories/log", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(&http.Cookie{Name: testCookieName, Value: testSessionToken})

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	// Form submissions successfully redirect back to the app instead of purely returning JSON
	if rr.Code != http.StatusSeeOther {
		t.Fatalf("expected status 303 Redirect for form, got %d", rr.Code)
	}
	if !store.logAdded {
		t.Fatalf("expected log to be added to store")
	}
}

func TestCaloriesUpdateTargetAPI_JSON(t *testing.T) {
	store := &mockCaloriesStore{}
	router := setupTestRouter(store)

	payload := map[string]any{
		"targetKcal":    2400,
		"targetProtein": 180,
		"targetCarbs":   220,
		"targetFat":     70,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal payload: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/calories/target", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{Name: testCookieName, Value: testSessionToken})

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: body: %s", rr.Code, rr.Body.String())
	}
}

func TestCaloriesDeleteLogAPI_Query(t *testing.T) {
	store := &mockCaloriesStore{}
	router := setupTestRouter(store)

	req := httptest.NewRequest(http.MethodDelete, "/api/calories/log?id=log-123", nil)
	req.AddCookie(&http.Cookie{Name: testCookieName, Value: testSessionToken})

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d: body: %s", rr.Code, rr.Body.String())
	}
}
