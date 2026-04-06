package authcontext

import "context"

type key string

const (
	userIDKey    key = "userID"
	sessionIDKey key = "sessionID"
)

func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

func WithSessionID(ctx context.Context, sessionID string) context.Context {
	return context.WithValue(ctx, sessionIDKey, sessionID)
}

func GetUserID(ctx context.Context) (string, bool) {
	val, ok := ctx.Value(userIDKey).(string)
	return val, ok
}

func GetSessionID(ctx context.Context) (string, bool) {
	val, ok := ctx.Value(sessionIDKey).(string)
	return val, ok
}
