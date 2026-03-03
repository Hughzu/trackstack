package httptransport

import (
	"context"
	"net/http"
	"strings"
)

type contextKey string

const (
	authUserIDKey    contextKey = "auth.user_id"
	authSessionIDKey contextKey = "auth.session_id"
)

func withAuthContext(ctx context.Context, userID string, sessionID string) context.Context {
	ctx = context.WithValue(ctx, authUserIDKey, userID)
	ctx = context.WithValue(ctx, authSessionIDKey, sessionID)
	return ctx
}

func authUserIDFromContext(ctx context.Context) string {
	value, _ := ctx.Value(authUserIDKey).(string)
	return value
}

func requireAuthUserID(w http.ResponseWriter, r *http.Request) (string, bool) {
	userID := strings.TrimSpace(authUserIDFromContext(r.Context()))
	if userID == "" {
		writeJSON(w, http.StatusUnauthorized, errorResponse{Error: "Unauthorized"})
		return "", false
	}

	return userID, true
}
