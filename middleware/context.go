package middleware

import (
	"context"
)

type middlewareContextKey string

const (
	userIDContextKey middlewareContextKey = "user_id_key"
)

func withUserIDString(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDContextKey, userID)
}

func UserIDStringFromContext(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(userIDContextKey).(string)
	if len(v) == 0 {
		ok = false
	}
	return v, ok
}
