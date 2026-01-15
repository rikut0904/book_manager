package authctx

import "context"

type AuthInfo struct {
	UserID string
	Email  string
	Name   string
}

type contextKey string

const authInfoKey contextKey = "authInfo"

func WithAuthInfo(ctx context.Context, info AuthInfo) context.Context {
	return context.WithValue(ctx, authInfoKey, info)
}

func AuthInfoFromContext(ctx context.Context) (AuthInfo, bool) {
	value := ctx.Value(authInfoKey)
	info, ok := value.(AuthInfo)
	return info, ok
}

func UserIDFromContext(ctx context.Context) string {
	if info, ok := AuthInfoFromContext(ctx); ok {
		return info.UserID
	}
	return ""
}
