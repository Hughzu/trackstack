package httptransport_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
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

const testPasswordHash = "$scrypt$N=16384,r=8,p=1$7skQJXQJzB4/K6qW2g6t2A==$MdanxXeVurzvktH6ZpHuqYg6pFVhq89zuAYxug0Tp6VpE4UJ7jfQ9XdDPaZ9Mor/6HYk/m8AKJ3MTXTRVrr19g=="

type authUsersStore struct {
	updated   bool
	updatedCh chan struct{}
}

func (m *authUsersStore) FindByEmail(ctx context.Context, email string) (users.User, error) {
	if email != "test@test.be" {
		return users.User{}, users.ErrNotFound
	}
	return users.User{ID: "test-user", Email: email, PasswordHash: testPasswordHash}, nil
}

func (m *authUsersStore) UpdateLastLogin(ctx context.Context, userID string, lastLoginAt string) error {
	m.updated = true
	if m.updatedCh != nil {
		select {
		case <-m.updatedCh:
		default:
			close(m.updatedCh)
		}
	}
	return nil
}

type authSessionStore struct {
	inserted bool
	revoked  bool
}

func (m *authSessionStore) FindSessionByID(ctx context.Context, id string) (auth.Session, error) {
	return auth.Session{}, auth.ErrUnauthorized
}

func (m *authSessionStore) InsertSession(ctx context.Context, session auth.Session) error {
	m.inserted = true
	return nil
}

func (m *authSessionStore) TouchSession(ctx context.Context, id string, lastSeenAt string, expiresAt string) error {
	return nil
}

func (m *authSessionStore) RotateOutSession(ctx context.Context, id string, revokedAt string, expiresAt string, rotatedAt string) error {
	return nil
}

func (m *authSessionStore) RevokeSession(ctx context.Context, id string, revokedAt string) error {
	m.revoked = true
	return nil
}

func setupAuthTestRouter(sessionStore auth.SessionStore, usersStore *authUsersStore) *chi.Mux {
	authService := auth.NewService(sessionStore, testAuthConfig())
	usersService := users.NewService(usersStore)

	handlers := httptransport.NewHandlers(httptransport.Deps{
		AuthService:        authService,
		UsersService:       usersService,
		CaloriesService:    calories.NewService(nil),
		HeatService:        heat.NewService(nil),
		ExpensesService:    expenses.NewService(nil),
		AuthCookieName:     testCookieName,
		AuthCookieSecure:   false,
		AuthCookieSameSite: "lax",
	})

	return httptransport.NewRouter(handlers, "*").(*chi.Mux)
}

func TestAuthLoginAPI_JSON(t *testing.T) {
	sessionStore := &authSessionStore{}
	usersStore := &authUsersStore{updatedCh: make(chan struct{})}
	router := setupAuthTestRouter(sessionStore, usersStore)

	// TODO, remove the hardcoded password, get this from .env, rotate password ...
	body, err := json.Marshal(map[string]any{
		"email":    "test@test.be",
		"password": "Test123*",
	})
	if err != nil {
		t.Fatalf("marshal payload: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d: body: %s", rr.Code, rr.Body.String())
	}
	if !sessionStore.inserted {
		t.Fatalf("expected session to be created")
	}
	select {
	case <-usersStore.updatedCh:
	case <-time.After(time.Second):
		t.Fatalf("timed out waiting for last login update")
	}
	if !usersStore.updated {
		t.Fatalf("expected last login to be updated")
	}
	if len(rr.Result().Cookies()) == 0 {
		t.Fatalf("expected auth cookie to be set")
	}
}

func TestAuthLogoutAPI_JSON(t *testing.T) {
	sessionStore := &authSessionStore{}
	usersStore := &authUsersStore{}
	router := setupAuthTestRouter(sessionStore, usersStore)

	req := httptest.NewRequest(http.MethodPost, "/api/auth/logout", nil)
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{Name: testCookieName, Value: testSessionToken})
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d: body: %s", rr.Code, rr.Body.String())
	}
	if !sessionStore.revoked {
		t.Fatalf("expected session to be revoked")
	}
	if len(rr.Result().Cookies()) == 0 {
		t.Fatalf("expected auth cookie clearing response")
	}
}

func TestAuthSessionAPI_JSON(t *testing.T) {
	sessionStore := &testAuthStore{}
	usersStore := &authUsersStore{}
	router := setupAuthTestRouter(sessionStore, usersStore)

	req := httptest.NewRequest(http.MethodGet, "/api/auth/session", nil)
	req.AddCookie(&http.Cookie{Name: testCookieName, Value: testSessionToken})
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: body: %s", rr.Code, rr.Body.String())
	}
	if !bytes.Contains(rr.Body.Bytes(), []byte(`"userId":"test-user"`)) {
		t.Fatalf("expected user id in body, got %s", rr.Body.String())
	}
}
