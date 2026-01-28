package middleware

import (
	"net/http"
	"strings"

	"book_manager/backend/internal/authctx"
	"book_manager/backend/internal/firebaseauth"
	"book_manager/backend/internal/handler"
	"book_manager/backend/internal/users"
)

type FirebaseAuthMiddleware struct {
	verifier     *firebaseauth.Verifier
	usersService *users.Service
}

func NewFirebaseAuthMiddleware(verifier *firebaseauth.Verifier, usersService *users.Service) *FirebaseAuthMiddleware {
	return &FirebaseAuthMiddleware{
		verifier:     verifier,
		usersService: usersService,
	}
}

func (m *FirebaseAuthMiddleware) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/healthz" || r.Method == http.MethodOptions || strings.HasPrefix(r.URL.Path, "/auth/") {
			next.ServeHTTP(w, r)
			return
		}
		token := bearerToken(r.Header.Get("Authorization"))
		if token == "" {
			handler.Unauthorized(w)
			return
		}
		info, err := m.verifier.VerifyIDToken(r.Context(), token)
		if err != nil {
			handler.Unauthorized(w)
			return
		}
		if m.usersService != nil {
			if _, err := m.usersService.Ensure(info.UserID, info.Email, info.Name); err != nil {
				handler.InternalError(w)
				return
			}
		}
		if !info.EmailVerified {
			handler.EmailNotVerified(w)
			return
		}
		ctx := authctx.WithAuthInfo(r.Context(), authctx.AuthInfo{
			UserID: info.UserID,
			Email:  info.Email,
			Name:   info.Name,
			EmailVerified: info.EmailVerified,
		})
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func bearerToken(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	parts := strings.SplitN(value, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return ""
	}
	return strings.TrimSpace(parts[1])
}
