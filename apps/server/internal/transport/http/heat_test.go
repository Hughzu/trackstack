package httptransport_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
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

type mockHeatStore struct {
	created bool
	deleted string
	recent  []heat.Refill
}

func (m *mockHeatStore) ListByRange(ctx context.Context, userID string, from string, to string) ([]heat.Refill, error) {
	return nil, nil
}

func (m *mockHeatStore) ListRecent(ctx context.Context, userID string, limit int, offset int) ([]heat.Refill, error) {
	if m.recent == nil {
		return []heat.Refill{}, nil
	}
	return m.recent, nil
}

func (m *mockHeatStore) Create(ctx context.Context, refill heat.Refill) error {
	m.created = true
	return nil
}

func (m *mockHeatStore) GetLatest(ctx context.Context, userID string) (*heat.Refill, error) {
	return nil, nil
}

func (m *mockHeatStore) GetSumByRange(ctx context.Context, userID string, fromDate string, toDate string) (int, error) {
	return 0, nil
}

func (m *mockHeatStore) Delete(ctx context.Context, userID string, id string) (bool, error) {
	m.deleted = id
	return true, nil
}

func setupHeatTestRouter(hStore *mockHeatStore) *chi.Mux {
	authService := auth.NewService(&testAuthStore{}, testAuthConfig())
	usersService := users.NewService(&testUsersStore{})
	heatService := heat.NewService(hStore)

	handlers := httptransport.NewHandlers(httptransport.Deps{
		AuthService:        authService,
		UsersService:       usersService,
		CaloriesService:    calories.NewService(nil),
		HeatService:        heatService,
		ExpensesService:    expenses.NewService(nil),
		AuthCookieName:     testCookieName,
		AuthCookieSecure:   false,
		AuthCookieSameSite: "lax",
	})

	return httptransport.NewRouter(handlers, "*").(*chi.Mux)
}

func TestHeatCreateRefillAPI_JSON(t *testing.T) {
	store := &mockHeatStore{}
	router := setupHeatTestRouter(store)

	payload := map[string]any{
		"date":     "2026-03-10",
		"weightKg": 30,
		"bags":     2,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal payload: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/heat/refills", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{Name: testCookieName, Value: testSessionToken})

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d: body: %s", rr.Code, rr.Body.String())
	}
	if !store.created {
		t.Fatalf("expected refill to be created")
	}
}

func TestHeatDeleteRefillAPI_Query(t *testing.T) {
	store := &mockHeatStore{}
	router := setupHeatTestRouter(store)

	req := httptest.NewRequest(http.MethodDelete, "/api/heat/refills?id=refill-123", nil)
	req.AddCookie(&http.Cookie{Name: testCookieName, Value: testSessionToken})

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d: body: %s", rr.Code, rr.Body.String())
	}
	if store.deleted != "refill-123" {
		t.Fatalf("expected refill id to be deleted, got %q", store.deleted)
	}
}

func TestHeatDeleteRefillAPI_BodyRejected(t *testing.T) {
	store := &mockHeatStore{}
	router := setupHeatTestRouter(store)

	req := httptest.NewRequest(http.MethodDelete, "/api/heat/refills", strings.NewReader(`{"id":"refill-123"}`))
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{Name: testCookieName, Value: testSessionToken})

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400 when id missing from query, got %d: body: %s", rr.Code, rr.Body.String())
	}
	if store.deleted != "" {
		t.Fatalf("expected delete not to be called, got %q", store.deleted)
	}
}
