package authcontext

import "context"

type key string

const userIDKey key = "userID"

func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

func GetUserID(ctx context.Context) (string, bool) {
	val, ok := ctx.Value(userIDKey).(string)
	return val, ok
}
