package router

import (
	"net/http"

	"book_manager/backend/internal/handler"
	"book_manager/backend/internal/repository"
)

func New(
	h *handler.Handler,
	auditRepo repository.AuditLogRepository,
	allowedOrigins string,
) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", h.Health)

	mux.HandleFunc("/auth/signup", h.AuthSignup)
	mux.HandleFunc("/auth/login", h.AuthLogin)
	mux.HandleFunc("/auth/refresh", h.AuthRefresh)
	mux.HandleFunc("/auth/logout", h.AuthLogout)

	mux.HandleFunc("/isbn/lookup", h.IsbnLookup)

	mux.HandleFunc("/books", h.Books)
	mux.HandleFunc("/books/", h.BookByID)

	mux.HandleFunc("/user-books", h.UserBooks)
	mux.HandleFunc("/user-books/", h.UserBooksByID)
	mux.HandleFunc("/user-series/override", h.UserSeriesOverride)

	mux.HandleFunc("/favorites", h.Favorites)
	mux.HandleFunc("/favorites/", h.FavoritesByID)

	mux.HandleFunc("/next-to-buy", h.NextToBuy)
	mux.HandleFunc("/next-to-buy/manual", h.NextToBuyManual)
	mux.HandleFunc("/next-to-buy/manual/", h.NextToBuyManualByID)

	mux.HandleFunc("/tags", h.Tags)
	mux.HandleFunc("/tags/", h.TagsByID)
	mux.HandleFunc("/book-tags", h.BookTags)

	mux.HandleFunc("/recommendations", h.Recommendations)
	mux.HandleFunc("/recommendations/", h.RecommendationsByID)

	mux.HandleFunc("/users", h.Users)
	mux.HandleFunc("/users/", h.UsersByID)
	mux.HandleFunc("/users/me", h.UsersMe)
	mux.HandleFunc("/users/me/settings", h.UsersMeSettings)

	mux.HandleFunc("/follows/", h.Follows)

	mux.HandleFunc("/book-reports", h.BookReports)

	mux.HandleFunc("/series", h.Series)
	mux.HandleFunc("/series/", h.SeriesByID)

	mux.HandleFunc("/admin/openai-keys", h.AdminOpenAIKeys)
	mux.HandleFunc("/admin/openai-keys/", h.AdminOpenAIKeysByID)
	mux.HandleFunc("/admin/openai-models", h.AdminOpenAIModels)

	handlerWithAudit := AuditMiddleware(auditRepo, mux)
	return CORSMiddleware(allowedOrigins, handlerWithAudit)
}
