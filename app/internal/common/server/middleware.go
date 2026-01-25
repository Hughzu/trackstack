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

// SessionMiddleware handles session cookie validation/creation
func SessionMiddleware(database *db.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			var user *db.User

			// 1. Check for session cookie
			cookie, err := r.Cookie("trackstack_session")
			if err == nil && cookie.Value != "" {
				// 2. Validate session from DB
				session, err := database.GetSessionByID(ctx, cookie.Value)
				if err == nil {
					// 3. Session is valid, load user
					user, err = database.GetUserByID(ctx, session.UserID)
					if err == nil {
						// Update last seen timestamp
						if err := database.UpdateUserLastSeen(ctx, user.ID); err != nil {
							slog.Error("failed to update last seen", "error", err, "user_id", user.ID)
						}
					} else {
						slog.Warn("session user not found", "user_id", session.UserID)
						user = nil
					}
				} else {
					if !errors.Is(err, sql.ErrNoRows) {
						slog.Warn("invalid session", "session_id", cookie.Value, "error", err)
					}
				}
			}

			// 4. If no valid user, create new user + session (auto-provision)
			if user == nil {
				user = db.NewUser()
				if err := database.CreateUser(ctx, user); err != nil {
					slog.Error("failed to create user", "error", err)
					http.Error(w, "Internal server error", http.StatusInternalServerError)
					return
				}

				session := db.NewSession(user.ID)
				if err := database.CreateSession(ctx, session); err != nil {
					slog.Error("failed to create session", "error", err)
					http.Error(w, "Internal server error", http.StatusInternalServerError)
					return
				}

				// 5. Set session cookie
				http.SetCookie(w, createSessionCookie(session.ID))

				slog.Info("created new user and session",
					"user_id", user.ID,
					"session_id", session.ID,
				)
			}

			// 6. Inject user into context
			ctx = context.WithValue(ctx, UserContextKey, user)

			// 7. Continue to next handler
			next.ServeHTTP(w, r.WithContext(ctx))
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
