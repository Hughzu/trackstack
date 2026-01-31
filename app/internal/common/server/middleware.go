package server

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"net/http"
	"os"

	"github.com/23St/trackstack/internal/common/db"
)

type contextKey string

const UserContextKey contextKey = "user"

// SessionMiddleware validates existing session cookies
// Note: User should already be set in context by IPWhitelistMiddleware
// This middleware validates the session and updates activity timestamps
func SessionMiddleware(database *db.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			// Check if user already in context (set by IP whitelist)
			if user, ok := GetUserFromContext(ctx); ok && user != nil {
				// User already validated by IP whitelist, just update last seen
				if err := database.UpdateUserLastSeen(ctx, user.ID); err != nil {
					slog.Error("failed to update last seen", "error", err, "user_id", user.ID)
				}
				next.ServeHTTP(w, r)
				return
			}

			// No user in context, try to validate session cookie
			cookie, err := r.Cookie("trackstack_session")
			if err == nil && cookie.Value != "" {
				session, err := database.GetSessionByID(ctx, cookie.Value)
				if err == nil {
					user, err := database.GetUserByID(ctx, session.UserID)
					if err == nil {
						// Valid session, update last seen and inject user
						if err := database.UpdateUserLastSeen(ctx, user.ID); err != nil {
							slog.Error("failed to update last seen", "error", err, "user_id", user.ID)
						}
						ctx = context.WithValue(ctx, UserContextKey, user)
						next.ServeHTTP(w, r.WithContext(ctx))
						return
					}
				} else {
					if !errors.Is(err, sql.ErrNoRows) {
						slog.Warn("invalid session", "session_id", cookie.Value, "error", err)
					}
				}
			}

			// No valid session - reject request
			// Note: IP whitelist should have caught this, but handle gracefully
			slog.Warn("no valid session or IP mapping", "path", r.URL.Path)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		})
	}
}

// createSessionCookie creates a session cookie with appropriate security flags
func createSessionCookie(sessionID string) *http.Cookie {
	isProduction := os.Getenv("ENV") == "production"

	return &http.Cookie{
		Name:     "trackstack_session",
		Value:    sessionID,
		Path:     "/",
		MaxAge:   2592000, // 30 days
		HttpOnly: true,
		Secure:   isProduction, // HTTPS only in production
		SameSite: http.SameSiteStrictMode,
	}
}

// GetUserFromContext retrieves user from request context
func GetUserFromContext(ctx context.Context) (*db.User, bool) {
	user, ok := ctx.Value(UserContextKey).(*db.User)
	return user, ok
}
