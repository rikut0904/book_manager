package handler

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"book_manager/backend/internal/auth"
	"book_manager/backend/internal/books"
	"book_manager/backend/internal/domain"
	"book_manager/backend/internal/follows"
	"book_manager/backend/internal/isbn"
	"book_manager/backend/internal/userbooks"
	"book_manager/backend/internal/users"
)

type Handler struct {
	auth      *auth.Service
	isbn      *isbn.Service
	books     *books.Service
	userBooks *userbooks.Service
	users     *users.Service
	follows   *follows.Service
}

func New(
	authService *auth.Service,
	isbnService *isbn.Service,
	bookService *books.Service,
	userBookService *userbooks.Service,
	usersService *users.Service,
	followsService *follows.Service,
) *Handler {
	return &Handler{
		auth:      authService,
		isbn:      isbnService,
		books:     bookService,
		userBooks: userBookService,
		users:     usersService,
		follows:   followsService,
	}
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"status": "ok",
	})
}

func (h *Handler) AuthSignup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		methodNotAllowed(w, http.MethodPost)
		return
	}
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Username string `json:"username"`
	}
	if err := decodeJSON(r, &req); err != nil {
		badRequest(w, "invalid json")
		return
	}
	if req.Email == "" || req.Password == "" || req.Username == "" {
		badRequest(w, "email, password, username are required")
		return
	}
	result, err := h.auth.Signup(req.Email, req.Password, req.Username)
	if err != nil {
		if errors.Is(err, auth.ErrUserExists) {
			conflict(w, "user already exists")
			return
		}
		internalError(w)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"accessToken":  result.AccessToken,
		"refreshToken": result.RefreshToken,
		"user": map[string]string{
			"id":       result.User.ID,
			"email":    result.User.Email,
			"username": result.User.Username,
		},
	})
}

func (h *Handler) AuthLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		methodNotAllowed(w, http.MethodPost)
		return
	}
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := decodeJSON(r, &req); err != nil {
		badRequest(w, "invalid json")
		return
	}
	if req.Email == "" || req.Password == "" {
		badRequest(w, "email and password are required")
		return
	}
	result, err := h.auth.Login(req.Email, req.Password)
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			unauthorized(w)
			return
		}
		internalError(w)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"accessToken":  result.AccessToken,
		"refreshToken": result.RefreshToken,
		"user": map[string]string{
			"id":       result.User.ID,
			"email":    result.User.Email,
			"username": result.User.Username,
		},
	})
}

func (h *Handler) AuthRefresh(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		methodNotAllowed(w, http.MethodPost)
		return
	}
	var req struct {
		RefreshToken string `json:"refreshToken"`
	}
	if err := decodeJSON(r, &req); err != nil {
		badRequest(w, "invalid json")
		return
	}
	if req.RefreshToken == "" {
		badRequest(w, "refreshToken is required")
		return
	}
	accessToken, refreshToken, err := h.auth.Refresh(req.RefreshToken)
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			unauthorized(w)
			return
		}
		internalError(w)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	})
}

func (h *Handler) AuthLogout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		methodNotAllowed(w, http.MethodPost)
		return
	}
	var req struct {
		RefreshToken string `json:"refreshToken"`
	}
	if err := decodeJSON(r, &req); err != nil {
		badRequest(w, "invalid json")
		return
	}
	if req.RefreshToken == "" {
		badRequest(w, "refreshToken is required")
		return
	}
	h.auth.Logout(req.RefreshToken)
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func (h *Handler) IsbnLookup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		methodNotAllowed(w, http.MethodGet)
		return
	}
	if r.URL.Query().Get("isbn") == "" {
		badRequest(w, "isbn is required")
		return
	}
	isbnValue := r.URL.Query().Get("isbn")
	book, err := h.isbn.Lookup(isbnValue)
	if err != nil {
		if errors.Is(err, isbn.ErrNotFound) {
			notFound(w)
			return
		}
		internalError(w)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"isbn13":        book.ISBN13,
		"title":         book.Title,
		"authors":       book.Authors,
		"publisher":     book.Publisher,
		"publishedDate": book.PublishedDate,
		"thumbnailUrl":  book.ThumbnailURL,
		"source":        book.Source,
	})
}

func (h *Handler) Books(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		items := h.books.List()
		writeJSON(w, http.StatusOK, map[string]any{
			"items": items,
		})
	case http.MethodPost:
		var req bookRequest
		if err := decodeJSON(r, &req); err != nil {
			badRequest(w, "invalid json")
			return
		}
		isbn13, err := validateBookRequest(req)
		if err != nil {
			badRequest(w, err.Error())
			return
		}
		book, err := h.books.Create(domain.Book{
			ISBN13:        isbn13,
			Title:         req.Title,
			Authors:       req.Authors,
			Publisher:     req.Publisher,
			PublishedDate: req.PublishedDate,
			ThumbnailURL:  req.ThumbnailURL,
			Source:        req.Source,
		})
		if err != nil {
			if errors.Is(err, books.ErrBookExists) {
				conflict(w, "book already exists")
				return
			}
			internalError(w)
			return
		}
		writeJSON(w, http.StatusOK, book)
	default:
		methodNotAllowed(w, http.MethodGet, http.MethodPost)
	}
}

func (h *Handler) BookByID(w http.ResponseWriter, r *http.Request) {
	if _, ok := pathID("/books/", r.URL.Path); !ok {
		notFound(w)
		return
	}
	if r.Method != http.MethodGet {
		methodNotAllowed(w, http.MethodGet)
		return
	}
	if id, ok := pathID("/books/", r.URL.Path); ok {
		if book, found := h.books.Get(id); found {
			writeJSON(w, http.StatusOK, book)
			return
		}
	}
	notFound(w)
}

func (h *Handler) UserBooks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		userID := r.URL.Query().Get("userId")
		if userID == "" {
			userID = "user_demo"
		}
		items := h.userBooks.ListByUser(userID)
		writeJSON(w, http.StatusOK, map[string]any{
			"items": items,
		})
	case http.MethodPost:
		var req struct {
			UserID     string `json:"userId"`
			BookID     string `json:"bookId"`
			Note       string `json:"note"`
			AcquiredAt string `json:"acquiredAt"`
		}
		if err := decodeJSON(r, &req); err != nil {
			badRequest(w, "invalid json")
			return
		}
		if strings.TrimSpace(req.BookID) == "" {
			badRequest(w, "bookId is required")
			return
		}
		userID := strings.TrimSpace(req.UserID)
		if userID == "" {
			userID = "user_demo"
		}
		if req.AcquiredAt != "" && !isISODate(req.AcquiredAt) {
			badRequest(w, "acquiredAt must be YYYY-MM-DD")
			return
		}
		item, err := h.userBooks.Create(userID, req.BookID, req.Note, req.AcquiredAt)
		if err != nil {
			if errors.Is(err, userbooks.ErrUserBookExists) {
				conflict(w, "user book already exists")
				return
			}
			internalError(w)
			return
		}
		writeJSON(w, http.StatusOK, item)
	default:
		methodNotAllowed(w, http.MethodGet, http.MethodPost)
	}
}

func (h *Handler) UserBooksByID(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPatch:
		id, ok := pathID("/user-books/", r.URL.Path)
		if !ok {
			notFound(w)
			return
		}
		var req struct {
			Note       *string `json:"note"`
			AcquiredAt *string `json:"acquiredAt"`
		}
		if err := decodeJSON(r, &req); err != nil {
			badRequest(w, "invalid json")
			return
		}
		if req.Note == nil && req.AcquiredAt == nil {
			badRequest(w, "note or acquiredAt is required")
			return
		}
		if req.AcquiredAt != nil && *req.AcquiredAt != "" && !isISODate(*req.AcquiredAt) {
			badRequest(w, "acquiredAt must be YYYY-MM-DD")
			return
		}
		item, ok := h.userBooks.Update(id, userbooks.UpdateInput{
			Note:       req.Note,
			AcquiredAt: req.AcquiredAt,
		})
		if !ok {
			notFound(w)
			return
		}
		writeJSON(w, http.StatusOK, item)
	case http.MethodDelete:
		id, ok := pathID("/user-books/", r.URL.Path)
		if !ok {
			notFound(w)
			return
		}
		if !h.userBooks.Delete(id) {
			notFound(w)
			return
		}
		writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
	default:
		methodNotAllowed(w, http.MethodPatch, http.MethodDelete)
	}
}

func (h *Handler) UserSeriesOverride(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		methodNotAllowed(w, http.MethodPatch)
		return
	}
	echoRequest(w, r)
}

func (h *Handler) Favorites(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		echoRequest(w, r)
	case http.MethodPost:
		echoRequest(w, r)
	default:
		methodNotAllowed(w, http.MethodGet, http.MethodPost)
	}
}

func (h *Handler) FavoritesByID(w http.ResponseWriter, r *http.Request) {
	if _, ok := pathID("/favorites/", r.URL.Path); !ok {
		notFound(w)
		return
	}
	if r.Method != http.MethodDelete {
		methodNotAllowed(w, http.MethodDelete)
		return
	}
	echoRequest(w, r)
}

func (h *Handler) NextToBuy(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		methodNotAllowed(w, http.MethodGet)
		return
	}
	echoRequest(w, r)
}

func (h *Handler) NextToBuyManual(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		methodNotAllowed(w, http.MethodPost)
		return
	}
	echoRequest(w, r)
}

func (h *Handler) NextToBuyManualByID(w http.ResponseWriter, r *http.Request) {
	if _, ok := pathID("/next-to-buy/manual/", r.URL.Path); !ok {
		notFound(w)
		return
	}
	switch r.Method {
	case http.MethodPatch:
		echoRequest(w, r)
	case http.MethodDelete:
		echoRequest(w, r)
	default:
		methodNotAllowed(w, http.MethodPatch, http.MethodDelete)
	}
}

func (h *Handler) Tags(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		echoRequest(w, r)
	case http.MethodPost:
		echoRequest(w, r)
	default:
		methodNotAllowed(w, http.MethodGet, http.MethodPost)
	}
}

func (h *Handler) TagsByID(w http.ResponseWriter, r *http.Request) {
	if _, ok := pathID("/tags/", r.URL.Path); !ok {
		notFound(w)
		return
	}
	if r.Method != http.MethodDelete {
		methodNotAllowed(w, http.MethodDelete)
		return
	}
	echoRequest(w, r)
}

func (h *Handler) BookTags(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		echoRequest(w, r)
	case http.MethodDelete:
		echoRequest(w, r)
	default:
		methodNotAllowed(w, http.MethodPost, http.MethodDelete)
	}
}

func (h *Handler) Recommendations(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		echoRequest(w, r)
	case http.MethodPost:
		echoRequest(w, r)
	default:
		methodNotAllowed(w, http.MethodGet, http.MethodPost)
	}
}

func (h *Handler) RecommendationsByID(w http.ResponseWriter, r *http.Request) {
	if _, ok := pathID("/recommendations/", r.URL.Path); !ok {
		notFound(w)
		return
	}
	if r.Method != http.MethodDelete {
		methodNotAllowed(w, http.MethodDelete)
		return
	}
	echoRequest(w, r)
}

func (h *Handler) Users(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		methodNotAllowed(w, http.MethodGet)
		return
	}
	query := r.URL.Query().Get("query")
	items := h.users.List(query)
	result := make([]map[string]string, 0, len(items))
	for _, user := range items {
		result = append(result, map[string]string{
			"id":       user.ID,
			"email":    user.Email,
			"username": user.Username,
		})
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"items": result,
	})
}

func (h *Handler) UsersByID(w http.ResponseWriter, r *http.Request) {
	if _, ok := pathID("/users/", r.URL.Path); !ok {
		notFound(w)
		return
	}
	if r.Method != http.MethodGet {
		methodNotAllowed(w, http.MethodGet)
		return
	}
	userID, _ := pathID("/users/", r.URL.Path)
	user, ok := h.users.Get(userID)
	if !ok {
		notFound(w)
		return
	}
	settings := h.users.GetSettings(userID)
	ownedCount := len(h.userBooks.ListByUser(userID))
	writeJSON(w, http.StatusOK, map[string]any{
		"user": map[string]string{
			"id":       user.ID,
			"email":    user.Email,
			"username": user.Username,
		},
		"settings": map[string]string{
			"visibility": settings.Visibility,
		},
		"stats": map[string]int{
			"ownedCount":  ownedCount,
			"seriesCount": 0,
			"followers":   h.follows.CountFollowers(userID),
			"following":   h.follows.CountFollowing(userID),
		},
		"recommendations": []string{},
	})
}

func (h *Handler) UsersMe(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPatch:
		var req struct {
			Username string `json:"username"`
		}
		if err := decodeJSON(r, &req); err != nil {
			badRequest(w, "invalid json")
			return
		}
		if strings.TrimSpace(req.Username) == "" {
			badRequest(w, "username is required")
			return
		}
		userID := userIDFromRequest(r)
		user, ok := h.users.UpdateUsername(userID, req.Username)
		if !ok {
			notFound(w)
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{
			"user": map[string]string{
				"id":       user.ID,
				"email":    user.Email,
				"username": user.Username,
			},
		})
	case http.MethodDelete:
		userID := userIDFromRequest(r)
		if !h.users.Delete(userID) {
			notFound(w)
			return
		}
		writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
	default:
		methodNotAllowed(w, http.MethodPatch, http.MethodDelete)
	}
}

func (h *Handler) UsersMeSettings(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		methodNotAllowed(w, http.MethodPatch)
		return
	}
	var req struct {
		Visibility string `json:"visibility"`
	}
	if err := decodeJSON(r, &req); err != nil {
		badRequest(w, "invalid json")
		return
	}
	if strings.TrimSpace(req.Visibility) == "" {
		badRequest(w, "visibility is required")
		return
	}
	userID := userIDFromRequest(r)
	settings, err := h.users.UpdateSettings(userID, req.Visibility)
	if err != nil {
		if errors.Is(err, users.ErrInvalidVisibility) {
			badRequest(w, "visibility must be public or followers")
			return
		}
		internalError(w)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"settings": map[string]string{
			"visibility": settings.Visibility,
		},
	})
}

func (h *Handler) Follows(w http.ResponseWriter, r *http.Request) {
	if _, ok := pathID("/follows/", r.URL.Path); !ok {
		notFound(w)
		return
	}
	switch r.Method {
	case http.MethodPost:
		followeeID, _ := pathID("/follows/", r.URL.Path)
		followerID := userIDFromRequest(r)
		h.follows.Follow(followerID, followeeID)
		writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
	case http.MethodDelete:
		followeeID, _ := pathID("/follows/", r.URL.Path)
		followerID := userIDFromRequest(r)
		h.follows.Unfollow(followerID, followeeID)
		writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
	default:
		methodNotAllowed(w, http.MethodPost, http.MethodDelete)
	}
}

func (h *Handler) BookReports(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		methodNotAllowed(w, http.MethodPost)
		return
	}
	echoRequest(w, r)
}

func pathID(prefix, path string) (string, bool) {
	if !strings.HasPrefix(path, prefix) {
		return "", false
	}
	id := strings.TrimPrefix(path, prefix)
	if id == "" || strings.Contains(id, "/") {
		return "", false
	}
	return id, true
}

func userIDFromRequest(r *http.Request) string {
	if value := strings.TrimSpace(r.Header.Get("X-User-Id")); value != "" {
		return value
	}
	if value := strings.TrimSpace(r.URL.Query().Get("userId")); value != "" {
		return value
	}
	return "user_demo"
}

func validateBookRequest(req bookRequest) (string, error) {
	if strings.TrimSpace(req.Title) == "" {
		return "", errors.New("title is required")
	}
	isbn13 := strings.TrimSpace(req.ISBN13)
	if isbn13 != "" {
		normalized := normalizeISBN13(isbn13)
		if len(normalized) != 13 {
			return "", errors.New("isbn13 must be 13 digits")
		}
		isbn13 = normalized
	}
	if req.Source != "" && req.Source != "manual" && req.Source != "google" {
		return "", errors.New("source must be manual or google")
	}
	if req.PublishedDate != "" && !isISODate(req.PublishedDate) {
		return "", errors.New("publishedDate must be YYYY-MM-DD")
	}
	return isbn13, nil
}

type bookRequest struct {
	ISBN13        string   `json:"isbn13"`
	Title         string   `json:"title"`
	Authors       []string `json:"authors"`
	Publisher     string   `json:"publisher"`
	PublishedDate string   `json:"publishedDate"`
	ThumbnailURL  string   `json:"thumbnailUrl"`
	Source        string   `json:"source"`
}

func normalizeISBN13(value string) string {
	builder := strings.Builder{}
	for _, r := range value {
		if r >= '0' && r <= '9' {
			builder.WriteRune(r)
		}
	}
	return builder.String()
}

func isISODate(value string) bool {
	if len(value) != 10 {
		return false
	}
	_, err := time.Parse("2006-01-02", value)
	return err == nil
}
