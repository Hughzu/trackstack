package httptransport_test

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/Hughzu/trackstack/apps/server/internal/modules/auth"
	"github.com/Hughzu/trackstack/apps/server/internal/modules/users"
)

const testSessionToken = "valid-token"
const testCookieName = "session"

func hashSessionToken(value string) string {
	hash := sha256.Sum256([]byte(value))
	return hex.EncodeToString(hash[:])
}

func testAuthConfig() auth.Config {
	return auth.Config{
		SessionAbsoluteSeconds:      86400,
		SessionIdleSeconds:          3600,
		SessionRotateAfterSeconds:   1800,
		SessionRotationGraceSeconds: 30,
		SessionTouchSeconds:         60,
	}
}

type testAuthStore struct {
	session          auth.Session
	touchCalls       int
	touchedLastSeen  string
	touchedExpiresAt string
}

func (m *testAuthStore) FindSessionByID(ctx context.Context, id string) (auth.Session, error) {
	if id == hashSessionToken(testSessionToken) {
		if m.session.ID != "" || m.session.UserID != "" {
			return m.session, nil
		}

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

func (m *testAuthStore) InsertSession(ctx context.Context, session auth.Session) error { return nil }

func (m *testAuthStore) TouchSession(ctx context.Context, id string, lastSeenAt string, expiresAt string) error {
	m.touchCalls++
	m.touchedLastSeen = lastSeenAt
	m.touchedExpiresAt = expiresAt
	return nil
}

func (m *testAuthStore) RotateOutSession(ctx context.Context, id string, revokedAt string, expiresAt string, rotatedAt string) error {
	return nil
}

func (m *testAuthStore) RevokeSession(ctx context.Context, id string, revokedAt string) error {
	return nil
}

type testUsersStore struct{}

func (m *testUsersStore) FindByEmail(ctx context.Context, email string) (users.User, error) {
	return users.User{}, nil
}

func (m *testUsersStore) UpdateLastLogin(ctx context.Context, userID string, lastLoginAt string) error {
	return nil
}
