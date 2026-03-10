package httptransport_test

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/Hughzu/trackstack/apps/server/internal/modules/auth"
	"github.com/Hughzu/trackstack/apps/server/internal/modules/calories"
	"github.com/Hughzu/trackstack/apps/server/internal/modules/expenses"
	"github.com/Hughzu/trackstack/apps/server/internal/modules/heat"
	"github.com/Hughzu/trackstack/apps/server/internal/modules/users"
	httptransport "github.com/Hughzu/trackstack/apps/server/internal/transport/http"
	"github.com/go-chi/chi/v5"
)

const validSessionToken = "valid-token"

func hashSessionToken(value string) string {
	hash := sha256.Sum256([]byte(value))
	return hex.EncodeToString(hash[:])
}

// mockAuthStore satisfies auth.SessionStore for testing
type mockAuthStore struct{}

func (m *mockAuthStore) FindSessionByID(ctx context.Context, id string) (auth.Session, error) {
	if id == hashSessionToken(validSessionToken) {
		now := time.Now().UTC()
		return auth.Session{
			UserID:            "test-user",
			ID:                id,
			ExpiresAt:         now.Add(24 * time.Hour).Format(time.RFC3339),
			AbsoluteExpiresAt: now.Add(48 * time.Hour).Format(time.RFC3339),
			RotatedAt:         now.Format(time.RFC3339),
			LastSeenAt:        now.Format(time.RFC3339),
		}, nil
	}
	return auth.Session{}, auth.ErrUnauthorized
}
func (m *mockAuthStore) InsertSession(ctx context.Context, session auth.Session) error { return nil }
func (m *mockAuthStore) TouchSession(ctx context.Context, id string, lastSeenAt string, expiresAt string) error {
	return nil
}
func (m *mockAuthStore) RotateOutSession(ctx context.Context, id string, revokedAt string, expiresAt string, rotatedAt string) error {
	return nil
}
func (m *mockAuthStore) RevokeSession(ctx context.Context, id string, revokedAt string) error {
	return nil
}

type mockUsersStore struct{}

func (m *mockUsersStore) FindByEmail(ctx context.Context, email string) (users.User, error) {
	return users.User{}, nil
}
func (m *mockUsersStore) UpdateLastLogin(ctx context.Context, userID string, lastLoginAt string) error {
	return nil
}

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
	cfg := auth.Config{
		SessionAbsoluteSeconds:      86400,
		SessionIdleSeconds:          3600,
		SessionRotateAfterSeconds:   1800,
		SessionRotationGraceSeconds: 30,
		SessionTouchSeconds:         60,
	}

	authService := auth.NewService(&mockAuthStore{}, cfg)
	usersService := users.NewService(&mockUsersStore{})
	calService := calories.NewService(cStore)

	handlers := httptransport.NewHandlers(httptransport.Deps{
		AuthService:        authService,
		UsersService:       usersService,
		CaloriesService:    calService,
		HeatService:        &heat.Service{},
		ExpensesService:    &expenses.Service{},
		AuthCookieName:     "session",
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
	req.AddCookie(&http.Cookie{Name: "session", Value: validSessionToken})

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
	req.AddCookie(&http.Cookie{Name: "session", Value: validSessionToken})

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
