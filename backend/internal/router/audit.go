package router

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"io"
	"net/http"
	"strings"
	"time"

	"book_manager/backend/internal/authctx"
	"book_manager/backend/internal/domain"
	"book_manager/backend/internal/repository"
)

func AuditMiddleware(repo repository.AuditLogRepository, next http.Handler) http.Handler {
	if repo == nil {
		return next
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(io.LimitReader(r.Body, 1<<20))
		r.Body = io.NopCloser(bytes.NewReader(body))

		action := methodToAction(r.Method)
		entity, entityID := pathToEntity(r.URL.Path)
		payload := map[string]any{
			"method": r.Method,
			"path":   r.URL.Path,
			"query":  r.URL.Query(),
			"body":   string(body),
		}

		_ = repo.Create(domain.AuditLog{
			ID:        newAuditID(),
			UserID:    userIDFromRequest(r),
			Action:    action,
			Entity:    entity,
			EntityID:  entityID,
			Payload:   payload,
			IP:        requestIP(r),
			UserAgent: r.UserAgent(),
			CreatedAt: time.Now(),
		})

		next.ServeHTTP(w, r)
	})
}

func methodToAction(method string) string {
	switch method {
	case http.MethodPost:
		return "create"
	case http.MethodPatch, http.MethodPut:
		return "update"
	case http.MethodDelete:
		return "delete"
	default:
		return "read"
	}
}

func pathToEntity(path string) (string, string) {
	trimmed := strings.Trim(path, "/")
	if trimmed == "" {
		return "root", ""
	}
	parts := strings.Split(trimmed, "/")
	entity := parts[0]
	entityID := ""
	if len(parts) > 1 {
		entityID = parts[1]
	}
	return entity, entityID
}

func userIDFromRequest(r *http.Request) string {
	return strings.TrimSpace(authctx.UserIDFromContext(r.Context()))
}

func requestIP(r *http.Request) string {
	if value := strings.TrimSpace(r.Header.Get("X-Forwarded-For")); value != "" {
		return strings.Split(value, ",")[0]
	}
	return r.RemoteAddr
}

func newAuditID() string {
	seed := make([]byte, 16)
	if _, err := rand.Read(seed); err != nil {
		return "audit_demo"
	}
	return base64.RawURLEncoding.EncodeToString(seed)
}
