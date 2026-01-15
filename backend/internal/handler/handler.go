package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"book_manager/backend/internal/ai"
	"book_manager/backend/internal/auth"
	"book_manager/backend/internal/books"
	"book_manager/backend/internal/domain"
	"book_manager/backend/internal/favorites"
	"book_manager/backend/internal/follows"
	"book_manager/backend/internal/isbn"
	"book_manager/backend/internal/nexttobuy"
	"book_manager/backend/internal/openaikeys"
	"book_manager/backend/internal/recommendations"
	"book_manager/backend/internal/reports"
	"book_manager/backend/internal/series"
	"book_manager/backend/internal/userbooks"
	"book_manager/backend/internal/users"
)

type Handler struct {
	auth               *auth.Service
	isbn               *isbn.Service
	books              *books.Service
	userBooks          *userbooks.Service
	users              *users.Service
	follows            *follows.Service
	favorites          *favorites.Service
	nextToBuy          *nexttobuy.Service
	recs               *recommendations.Service
	reports            *reports.Service
	series             *series.Service
	openAIKeys         *openaikeys.Service
	openAIAPIKey       string
	openAIDefaultModel string
	aiPrompt           string
	adminUserIDs       map[string]struct{}
}

func New(
	authService *auth.Service,
	isbnService *isbn.Service,
	bookService *books.Service,
	userBookService *userbooks.Service,
	usersService *users.Service,
	followsService *follows.Service,
	favoritesService *favorites.Service,
	nextToBuyService *nexttobuy.Service,
	recsService *recommendations.Service,
	reportsService *reports.Service,
	seriesService *series.Service,
	openAIKeyService *openaikeys.Service,
	openAIAPIKey string,
	openAIDefaultModel string,
	aiPrompt string,
	adminUserIDs []string,
) *Handler {
	adminMap := make(map[string]struct{}, len(adminUserIDs))
	for _, id := range adminUserIDs {
		if strings.TrimSpace(id) == "" {
			continue
		}
		adminMap[id] = struct{}{}
	}
	return &Handler{
		auth:               authService,
		isbn:               isbnService,
		books:              bookService,
		userBooks:          userBookService,
		users:              usersService,
		follows:            followsService,
		favorites:          favoritesService,
		nextToBuy:          nextToBuyService,
		recs:               recsService,
		reports:            reportsService,
		series:             seriesService,
		openAIKeys:         openAIKeyService,
		openAIAPIKey:       openAIAPIKey,
		openAIDefaultModel: openAIDefaultModel,
		aiPrompt:           aiPrompt,
		adminUserIDs:       adminMap,
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
	if !isValidEmail(req.Email) {
		badRequest(w, "email is invalid")
		return
	}
	if len(req.Password) < 8 {
		badRequest(w, "password must be at least 8 characters")
		return
	}
	if strings.TrimSpace(req.Username) == "" {
		badRequest(w, "username is required")
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
	if !isValidEmail(req.Email) {
		badRequest(w, "email is invalid")
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
	if strings.TrimSpace(req.RefreshToken) == "" {
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
	if strings.TrimSpace(req.RefreshToken) == "" {
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
	isbnValue := strings.TrimSpace(r.URL.Query().Get("isbn"))
	if isbnValue == "" {
		badRequest(w, "isbn is required")
		return
	}
	userID := userIDFromRequest(r)
	settings := h.users.GetSettings(userID)

	var book domain.Book
	seriesGuess := isbn.SeriesGuess{}
	seriesSource := "simple"
	if existing, ok := h.books.FindByISBN(isbnValue); ok {
		book = existing
		seriesGuess = isbn.InferSeries(book.Title, book.SeriesName)
	} else {
		fetched, guess, err := h.isbn.Lookup(isbnValue)
		if err != nil {
			if errors.Is(err, isbn.ErrNotFound) {
				notFound(w)
				return
			}
			log.Printf("isbn lookup error: %v", err)
			internalError(w)
			return
		}
		book = fetched
		book.OriginalTitle = fetched.Title
		if cleaned := isbn.NormalizeTitle(book.Title); cleaned != "" {
			book.Title = cleaned
		}
		seriesGuess = guess
		created, err := h.books.Create(book)
		if err != nil {
			if existing, ok := h.books.FindByISBN(isbnValue); ok {
				book = existing
			} else if book.ISBN13 != "" {
				if existing, ok := h.books.FindByISBN(book.ISBN13); ok {
					book = existing
				} else if errors.Is(err, books.ErrBookExists) {
					log.Printf("isbn lookup: book exists but not found by isbn=%s", book.ISBN13)
					internalError(w)
					return
				} else {
					log.Printf("isbn lookup: book create error: %v", err)
					internalError(w)
					return
				}
			} else if errors.Is(err, books.ErrBookExists) {
				log.Printf("isbn lookup: book exists but not found by isbn=%s", isbnValue)
				internalError(w)
				return
			} else {
				log.Printf("isbn lookup: book create error: %v", err)
				internalError(w)
				return
			}
		} else {
			book = created
		}
	}
	sharedKey, hasShared := h.openAIKeys.First()
	apiKey := h.openAIAPIKey
	if hasShared {
		apiKey = sharedKey.APIKey
	}
	if apiKey != "" {
		rawTitle := book.OriginalTitle
		if rawTitle == "" {
			rawTitle = book.Title
		}
		model := h.openAIDefaultModel
		if settings.OpenAIModel != "" {
			model = settings.OpenAIModel
		}
		client := ai.NewOpenAIClient(apiKey, model, h.aiPrompt)
		guess, err := client.GuessSeries(r.Context(), ai.SeriesInput{
			Title:         rawTitle,
			RawTitle:      rawTitle,
			Authors:       book.Authors,
			Publisher:     book.Publisher,
			PublishedDate: book.PublishedDate,
			ISBN13:        book.ISBN13,
			SeriesName:    book.SeriesName,
		})
		if err != nil {
			log.Printf("openai guess error: %v", err)
		} else if guess.IsSeries && strings.TrimSpace(guess.Name) != "" {
			seriesGuess = isbn.SeriesGuess{
				Name:         guess.Name,
				VolumeNumber: guess.VolumeNumber,
			}
			seriesSource = "openai"
		} else if !guess.IsSeries {
			seriesGuess = isbn.SeriesGuess{}
			seriesSource = "openai"
		}
	}
	if seriesGuess.Name != "" {
		normalizedName := isbn.NormalizeSeriesName(seriesGuess.Name)
		seriesGuess.Name = normalizedName
	}
	if seriesGuess.Name == "" {
		seriesGuess.VolumeNumber = 0
	}
	seriesID := ""
	if seriesGuess.Name != "" {
		if item, err := h.series.Ensure(seriesGuess.Name); err == nil {
			seriesID = item.ID
		}
	}
	if seriesGuess.Name != "" && book.SeriesName == "" {
		book.SeriesName = seriesGuess.Name
		_ = h.books.Update(book)
	}
	if seriesID != "" || seriesGuess.VolumeNumber > 0 {
		userItems := h.userBooks.ListByUser(userID)
		var userBookID string
		for _, item := range userItems {
			if item.BookID == book.ID {
				userBookID = item.ID
				break
			}
		}
		if userBookID == "" {
			if created, err := h.userBooks.Create(userID, book.ID, "", ""); err == nil {
				userBookID = created.ID
			}
		}
		if userBookID != "" {
			input := userbooks.UpdateInput{}
			if seriesID != "" {
				seriesIDCopy := seriesID
				input.SeriesID = &seriesIDCopy
			}
			if seriesGuess.VolumeNumber > 0 {
				volume := seriesGuess.VolumeNumber
				input.VolumeNumber = &volume
			}
			source := "auto"
			input.SeriesSource = &source
			if input.SeriesID != nil || input.VolumeNumber != nil {
				_, _ = h.userBooks.Update(userBookID, input)
			}
		}
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"id":            book.ID,
		"isbn13":        book.ISBN13,
		"title":         book.Title,
		"authors":       book.Authors,
		"publisher":     book.Publisher,
		"publishedDate": book.PublishedDate,
		"thumbnailUrl":  book.ThumbnailURL,
		"source":        book.Source,
		"seriesId":      seriesID,
		"seriesName":    seriesGuess.Name,
		"volumeNumber":  seriesGuess.VolumeNumber,
		"seriesSource":  seriesSource,
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
		if req.IsSeries && strings.TrimSpace(req.SeriesName) == "" {
			badRequest(w, "seriesName is required when isSeries is true")
			return
		}
		if req.VolumeNumber != nil && *req.VolumeNumber <= 0 {
			badRequest(w, "volumeNumber must be positive")
			return
		}
		seriesName := ""
		if req.IsSeries {
			seriesName = isbn.NormalizeSeriesName(req.SeriesName)
		}
		title := strings.TrimSpace(req.Title)
		if cleaned := isbn.NormalizeTitle(title); cleaned != "" {
			title = cleaned
		}
		book, err := h.books.Create(domain.Book{
			ISBN13:        isbn13,
			Title:         title,
			Authors:       req.Authors,
			Publisher:     req.Publisher,
			PublishedDate: req.PublishedDate,
			ThumbnailURL:  req.ThumbnailURL,
			Source:        req.Source,
			SeriesName:    seriesName,
		})
		if err != nil {
			if errors.Is(err, books.ErrBookExists) {
				conflict(w, "book already exists")
				return
			}
			internalError(w)
			return
		}
		if req.IsSeries && seriesName != "" {
			userID := userIDFromRequest(r)
			seriesID := ""
			if item, err := h.series.Ensure(seriesName); err == nil {
				seriesID = item.ID
			}
			userBookID := ""
			if created, err := h.userBooks.Create(userID, book.ID, "", ""); err == nil {
				userBookID = created.ID
			}
			if userBookID != "" {
				input := userbooks.UpdateInput{}
				if seriesID != "" {
					seriesIDCopy := seriesID
					input.SeriesID = &seriesIDCopy
				}
				if req.VolumeNumber != nil {
					volume := *req.VolumeNumber
					input.VolumeNumber = &volume
				}
				if seriesID != "" {
					source := "manual"
					input.SeriesSource = &source
				}
				if input.SeriesID != nil || input.VolumeNumber != nil {
					_, _ = h.userBooks.Update(userBookID, input)
				}
			}
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
	switch r.Method {
	case http.MethodGet:
		if id, ok := pathID("/books/", r.URL.Path); ok {
			if book, found := h.books.Get(id); found {
				writeJSON(w, http.StatusOK, book)
				return
			}
		}
		notFound(w)
	case http.MethodDelete:
		id, ok := pathID("/books/", r.URL.Path)
		if !ok {
			notFound(w)
			return
		}
		if !h.books.Delete(id) {
			notFound(w)
			return
		}
		writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
	default:
		methodNotAllowed(w, http.MethodGet, http.MethodDelete)
	}
}

func (h *Handler) UserBooks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		userID := strings.TrimSpace(r.URL.Query().Get("userId"))
		if userID == "" {
			userID = userIDFromRequest(r)
		}
		bookID := strings.TrimSpace(r.URL.Query().Get("bookId"))
		query := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("query")))
		seriesID := strings.TrimSpace(r.URL.Query().Get("series"))
		page := 1
		if value := strings.TrimSpace(r.URL.Query().Get("page")); value != "" {
			if parsed, err := strconv.Atoi(value); err == nil && parsed > 0 {
				page = parsed
			}
		}
		items := h.userBooks.ListByUser(userID)
		if bookID != "" {
			filtered := items[:0]
			for _, item := range items {
				if item.BookID == bookID {
					filtered = append(filtered, item)
				}
			}
			items = filtered
		}
		if seriesID != "" {
			filtered := items[:0]
			for _, item := range items {
				if item.SeriesID == seriesID {
					filtered = append(filtered, item)
				}
			}
			items = filtered
		}
		if query != "" {
			booksByID := make(map[string]domain.Book)
			for _, book := range h.books.List() {
				booksByID[book.ID] = book
			}
			filtered := items[:0]
			for _, item := range items {
				book, ok := booksByID[item.BookID]
				if !ok {
					continue
				}
				if strings.Contains(strings.ToLower(book.Title), query) ||
					strings.Contains(strings.ToLower(book.ISBN13), query) ||
					containsAuthor(book.Authors, query) {
					filtered = append(filtered, item)
				}
			}
			items = filtered
		}
		const pageSize = 20
		start := (page - 1) * pageSize
		if start < len(items) {
			end := start + pageSize
			if end > len(items) {
				end = len(items)
			}
			items = items[start:end]
		} else if len(items) > 0 {
			items = []domain.UserBook{}
		}
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
	var req struct {
		BookID       string `json:"bookId"`
		SeriesID     string `json:"seriesId"`
		VolumeNumber *int   `json:"volumeNumber"`
		UserID       string `json:"userId"`
		IsSeries     *bool  `json:"isSeries"`
	}
	if err := decodeJSON(r, &req); err != nil {
		badRequest(w, "invalid json")
		return
	}
	if strings.TrimSpace(req.BookID) == "" {
		badRequest(w, "bookId is required")
		return
	}
	isSeries := true
	if req.IsSeries != nil {
		isSeries = *req.IsSeries
	}
	if isSeries {
		if strings.TrimSpace(req.SeriesID) == "" {
			badRequest(w, "seriesId is required")
			return
		}
		if req.VolumeNumber == nil {
			badRequest(w, "volumeNumber is required")
			return
		}
		if *req.VolumeNumber <= 0 {
			badRequest(w, "volumeNumber must be positive")
			return
		}
	} else {
		req.SeriesID = ""
		zero := 0
		req.VolumeNumber = &zero
	}
	userID := strings.TrimSpace(req.UserID)
	if userID == "" {
		userID = userIDFromRequest(r)
	}
	items := h.userBooks.ListByUser(userID)
	for _, item := range items {
		if item.BookID != req.BookID {
			continue
		}
		seriesID := req.SeriesID
		volume := *req.VolumeNumber
		source := "manual"
		updated, ok := h.userBooks.Update(item.ID, userbooks.UpdateInput{
			SeriesID:     &seriesID,
			VolumeNumber: &volume,
			SeriesSource: &source,
		})
		if !ok {
			internalError(w)
			return
		}
		writeJSON(w, http.StatusOK, updated)
		return
	}
	created, err := h.userBooks.Create(userID, req.BookID, "", "")
	if err != nil {
		internalError(w)
		return
	}
	seriesID := req.SeriesID
	volume := *req.VolumeNumber
	source := "manual"
	updated, ok := h.userBooks.Update(created.ID, userbooks.UpdateInput{
		SeriesID:     &seriesID,
		VolumeNumber: &volume,
		SeriesSource: &source,
	})
	if !ok {
		internalError(w)
		return
	}
	writeJSON(w, http.StatusOK, updated)
}

func (h *Handler) Favorites(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		userID := userIDFromRequest(r)
		items := h.favorites.ListByUser(userID)
		writeJSON(w, http.StatusOK, map[string]any{
			"items": items,
		})
	case http.MethodPost:
		var req struct {
			Type     string `json:"type"`
			BookID   string `json:"bookId"`
			SeriesID string `json:"seriesId"`
		}
		if err := decodeJSON(r, &req); err != nil {
			badRequest(w, "invalid json")
			return
		}
		favoriteType := strings.TrimSpace(req.Type)
		if favoriteType != "book" && favoriteType != "series" {
			badRequest(w, "type must be book or series")
			return
		}
		if favoriteType == "book" && strings.TrimSpace(req.BookID) == "" {
			badRequest(w, "bookId is required")
			return
		}
		if favoriteType == "series" && strings.TrimSpace(req.SeriesID) == "" {
			badRequest(w, "seriesId is required")
			return
		}
		if favoriteType == "book" && strings.TrimSpace(req.SeriesID) != "" {
			badRequest(w, "seriesId must be empty for book favorite")
			return
		}
		if favoriteType == "series" && strings.TrimSpace(req.BookID) != "" {
			badRequest(w, "bookId must be empty for series favorite")
			return
		}
		userID := userIDFromRequest(r)
		item, err := h.favorites.Create(userID, favoriteType, req.BookID, req.SeriesID)
		if err != nil {
			if errors.Is(err, favorites.ErrFavoriteExists) {
				conflict(w, "favorite already exists")
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

func (h *Handler) FavoritesByID(w http.ResponseWriter, r *http.Request) {
	if _, ok := pathID("/favorites/", r.URL.Path); !ok {
		notFound(w)
		return
	}
	if r.Method != http.MethodDelete {
		methodNotAllowed(w, http.MethodDelete)
		return
	}
	id, _ := pathID("/favorites/", r.URL.Path)
	if !h.favorites.Delete(id) {
		notFound(w)
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func (h *Handler) NextToBuy(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		methodNotAllowed(w, http.MethodGet)
		return
	}
	userID := userIDFromRequest(r)
	manualItems := h.nextToBuy.ListByUser(userID)
	items := make([]map[string]any, 0, len(manualItems))
	for _, item := range manualItems {
		items = append(items, map[string]any{
			"id":           item.ID,
			"title":        item.Title,
			"seriesName":   item.SeriesName,
			"volumeNumber": item.VolumeNumber,
			"note":         item.Note,
			"source":       "manual",
		})
	}
	seriesMap := make(map[string]string)
	for _, series := range h.series.List() {
		seriesMap[series.ID] = series.Name
	}
	maxBySeries := make(map[string]int)
	for _, item := range h.userBooks.ListByUser(userID) {
		if item.SeriesID == "" || item.VolumeNumber <= 0 {
			continue
		}
		if item.VolumeNumber > maxBySeries[item.SeriesID] {
			maxBySeries[item.SeriesID] = item.VolumeNumber
		}
	}
	for _, fav := range h.favorites.ListByUser(userID) {
		if fav.Type != "series" || fav.SeriesID == "" {
			continue
		}
		name := seriesMap[fav.SeriesID]
		nextVolume := maxBySeries[fav.SeriesID] + 1
		if nextVolume == 1 {
			nextVolume = 1
		}
		items = append(items, map[string]any{
			"id":           "auto:" + fav.SeriesID,
			"title":        name,
			"seriesName":   name,
			"volumeNumber": nextVolume,
			"note":         "お気に入りから自動提案",
			"source":       "auto",
		})
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"items": items,
	})
}

func (h *Handler) NextToBuyManual(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		methodNotAllowed(w, http.MethodPost)
		return
	}
	var req struct {
		Title        string `json:"title"`
		SeriesName   string `json:"seriesName"`
		VolumeNumber *int   `json:"volumeNumber"`
		Note         string `json:"note"`
	}
	if err := decodeJSON(r, &req); err != nil {
		badRequest(w, "invalid json")
		return
	}
	if strings.TrimSpace(req.Title) == "" {
		badRequest(w, "title is required")
		return
	}
	if req.VolumeNumber != nil && *req.VolumeNumber < 0 {
		badRequest(w, "volumeNumber must be 0 or positive")
		return
	}
	userID := userIDFromRequest(r)
	volume := 0
	if req.VolumeNumber != nil {
		volume = *req.VolumeNumber
	}
	item, err := h.nextToBuy.Create(userID, req.Title, req.SeriesName, volume, req.Note)
	if err != nil {
		internalError(w)
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (h *Handler) NextToBuyManualByID(w http.ResponseWriter, r *http.Request) {
	if _, ok := pathID("/next-to-buy/manual/", r.URL.Path); !ok {
		notFound(w)
		return
	}
	switch r.Method {
	case http.MethodPatch:
		id, ok := pathID("/next-to-buy/manual/", r.URL.Path)
		if !ok {
			notFound(w)
			return
		}
		var req struct {
			Title        *string `json:"title"`
			SeriesName   *string `json:"seriesName"`
			VolumeNumber *int    `json:"volumeNumber"`
			Note         *string `json:"note"`
		}
		if err := decodeJSON(r, &req); err != nil {
			badRequest(w, "invalid json")
			return
		}
		if req.Title == nil && req.SeriesName == nil && req.VolumeNumber == nil && req.Note == nil {
			badRequest(w, "no fields to update")
			return
		}
		if req.VolumeNumber != nil && *req.VolumeNumber < 0 {
			badRequest(w, "volumeNumber must be 0 or positive")
			return
		}
		item, ok := h.nextToBuy.Update(id, nexttobuy.UpdateInput{
			Title:        req.Title,
			SeriesName:   req.SeriesName,
			VolumeNumber: req.VolumeNumber,
			Note:         req.Note,
		})
		if !ok {
			notFound(w)
			return
		}
		writeJSON(w, http.StatusOK, item)
	case http.MethodDelete:
		id, ok := pathID("/next-to-buy/manual/", r.URL.Path)
		if !ok {
			notFound(w)
			return
		}
		if !h.nextToBuy.Delete(id) {
			notFound(w)
			return
		}
		writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
	default:
		methodNotAllowed(w, http.MethodPatch, http.MethodDelete)
	}
}

func (h *Handler) Recommendations(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		items := h.recs.List()
		writeJSON(w, http.StatusOK, map[string]any{
			"items": items,
		})
	case http.MethodPost:
		var req struct {
			BookID  string `json:"bookId"`
			Comment string `json:"comment"`
		}
		if err := decodeJSON(r, &req); err != nil {
			badRequest(w, "invalid json")
			return
		}
		if strings.TrimSpace(req.BookID) == "" {
			badRequest(w, "bookId is required")
			return
		}
		if strings.TrimSpace(req.Comment) == "" {
			badRequest(w, "comment is required")
			return
		}
		userID := userIDFromRequest(r)
		item, err := h.recs.Create(userID, req.BookID, req.Comment)
		if err != nil {
			internalError(w)
			return
		}
		writeJSON(w, http.StatusOK, item)
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
	id, _ := pathID("/recommendations/", r.URL.Path)
	if !h.recs.Delete(id) {
		notFound(w)
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
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
	seriesCount := len(h.series.List())
	writeJSON(w, http.StatusOK, map[string]any{
		"user": map[string]string{
			"id":       user.ID,
			"email":    user.Email,
			"username": user.Username,
		},
		"isAdmin": h.isAdminUser(user.ID),
		"settings": map[string]any{
			"visibility":    settings.Visibility,
			"openaiEnabled": settings.OpenAIEnabled,
			"openaiModel":   settings.OpenAIModel,
			"openaiHasKey":  settings.OpenAIAPIKey != "",
		},
		"stats": map[string]int{
			"ownedCount":  ownedCount,
			"seriesCount": seriesCount,
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
		for _, item := range h.recs.List() {
			if item.UserID == userID {
				h.recs.Delete(item.ID)
			}
		}
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
		Visibility    string `json:"visibility"`
		OpenAIEnabled *bool  `json:"openaiEnabled"`
		OpenAIModel   string `json:"openaiModel"`
		OpenAIAPIKey  string `json:"openaiApiKey"`
	}
	if err := decodeJSON(r, &req); err != nil {
		badRequest(w, "invalid json")
		return
	}
	if strings.TrimSpace(req.Visibility) == "" {
		req.Visibility = ""
	}
	userID := userIDFromRequest(r)
	if req.OpenAIEnabled != nil || req.OpenAIAPIKey != "" {
		badRequest(w, "use shared openai keys")
		return
	}
	if req.OpenAIModel != "" && !h.isAdminUser(userID) {
		forbidden(w, "admin only")
		return
	}
	settings, err := h.users.UpdateSettings(userID, req.Visibility, req.OpenAIEnabled, req.OpenAIModel, req.OpenAIAPIKey)
	if err != nil {
		if errors.Is(err, users.ErrInvalidVisibility) {
			badRequest(w, "visibility must be public or followers")
			return
		}
		internalError(w)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"settings": map[string]any{
			"visibility":    settings.Visibility,
			"openaiEnabled": settings.OpenAIEnabled,
			"openaiModel":   settings.OpenAIModel,
			"openaiHasKey":  settings.OpenAIAPIKey != "",
		},
	})
}

func (h *Handler) isAdminUser(userID string) bool {
	if userID == "" {
		return false
	}
	_, ok := h.adminUserIDs[userID]
	return ok
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
		if followeeID == followerID {
			badRequest(w, "cannot follow yourself")
			return
		}
		h.follows.Follow(followerID, followeeID)
		writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
	case http.MethodDelete:
		followeeID, _ := pathID("/follows/", r.URL.Path)
		followerID := userIDFromRequest(r)
		if followeeID == followerID {
			badRequest(w, "cannot unfollow yourself")
			return
		}
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
	var req struct {
		BookID     string `json:"bookId"`
		Suggestion string `json:"suggestion"`
		Note       string `json:"note"`
	}
	if err := decodeJSON(r, &req); err != nil {
		badRequest(w, "invalid json")
		return
	}
	if strings.TrimSpace(req.BookID) == "" || strings.TrimSpace(req.Suggestion) == "" {
		badRequest(w, "bookId and suggestion are required")
		return
	}
	book, ok := h.books.Get(req.BookID)
	if !ok {
		notFoundWithMessage(w, "book not found")
		return
	}
	h.reports.SendBookReport(reports.BookReport{
		BookID:     req.BookID,
		Suggestion: req.Suggestion,
		Note:       req.Note,
		Book:       book,
	})
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func (h *Handler) Series(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		items := h.series.List()
		writeJSON(w, http.StatusOK, map[string]any{
			"items": items,
		})
	case http.MethodPost:
		var req struct {
			Name string `json:"name"`
		}
		if err := decodeJSON(r, &req); err != nil {
			badRequest(w, "invalid json")
			return
		}
		if strings.TrimSpace(req.Name) == "" {
			badRequest(w, "name is required")
			return
		}
		item, err := h.series.Create(req.Name)
		if err != nil {
			if errors.Is(err, series.ErrSeriesExists) {
				conflict(w, "series already exists")
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

func (h *Handler) SeriesByID(w http.ResponseWriter, r *http.Request) {
	if _, ok := pathID("/series/", r.URL.Path); !ok {
		notFound(w)
		return
	}
	if r.Method != http.MethodDelete {
		methodNotAllowed(w, http.MethodDelete)
		return
	}
	id, _ := pathID("/series/", r.URL.Path)
	if len(h.userBooks.ListBySeriesID(id)) > 0 {
		conflict(w, "series is referenced by user books")
		return
	}
	if len(h.favorites.ListBySeriesID(id)) > 0 {
		conflict(w, "series is referenced by favorites")
		return
	}
	if !h.series.Delete(id) {
		notFound(w)
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func (h *Handler) AdminOpenAIKeys(w http.ResponseWriter, r *http.Request) {
	if !h.isAdminUser(userIDFromRequest(r)) {
		forbidden(w, "admin only")
		return
	}
	switch r.Method {
	case http.MethodGet:
		items := h.openAIKeys.List()
		out := make([]map[string]any, 0, len(items))
		for _, item := range items {
			out = append(out, map[string]any{
				"id":        item.ID,
				"name":      item.Name,
				"maskedKey": openaikeys.MaskKey(item.APIKey),
				"createdAt": item.CreatedAt,
			})
		}
		writeJSON(w, http.StatusOK, map[string]any{"items": out})
	case http.MethodPost:
		var req struct {
			Name   string `json:"name"`
			APIKey string `json:"apiKey"`
		}
		if err := decodeJSON(r, &req); err != nil {
			badRequest(w, "invalid json")
			return
		}
		if strings.TrimSpace(req.Name) == "" || strings.TrimSpace(req.APIKey) == "" {
			badRequest(w, "name and apiKey are required")
			return
		}
		item := h.openAIKeys.Create(req.Name, req.APIKey)
		writeJSON(w, http.StatusOK, map[string]any{
			"id":        item.ID,
			"name":      item.Name,
			"maskedKey": openaikeys.MaskKey(item.APIKey),
			"createdAt": item.CreatedAt,
		})
	default:
		methodNotAllowed(w, http.MethodGet, http.MethodPost)
	}
}

func (h *Handler) AdminOpenAIKeysByID(w http.ResponseWriter, r *http.Request) {
	if !h.isAdminUser(userIDFromRequest(r)) {
		forbidden(w, "admin only")
		return
	}
	if _, ok := pathID("/admin/openai-keys/", r.URL.Path); !ok {
		notFound(w)
		return
	}
	if r.Method != http.MethodDelete {
		methodNotAllowed(w, http.MethodDelete)
		return
	}
	id, _ := pathID("/admin/openai-keys/", r.URL.Path)
	if !h.openAIKeys.Delete(id) {
		notFound(w)
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func (h *Handler) AdminOpenAIModels(w http.ResponseWriter, r *http.Request) {
	if !h.isAdminUser(userIDFromRequest(r)) {
		forbidden(w, "admin only")
		return
	}
	if r.Method != http.MethodGet {
		methodNotAllowed(w, http.MethodGet)
		return
	}
	sharedKey, ok := h.openAIKeys.First()
	apiKey := h.openAIAPIKey
	if ok && strings.TrimSpace(sharedKey.APIKey) != "" {
		apiKey = sharedKey.APIKey
	}
	if strings.TrimSpace(apiKey) == "" {
		badRequest(w, "shared api key is required")
		return
	}
	models, err := fetchOpenAIModels(apiKey)
	if err != nil {
		internalError(w)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": models})
}

type openAIModelsResponse struct {
	Data []struct {
		ID string `json:"id"`
	} `json:"data"`
}

func fetchOpenAIModels(apiKey string) ([]string, error) {
	req, err := http.NewRequest(http.MethodGet, "https://api.openai.com/v1/models", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("openai models fetch failed")
	}
	var payload openAIModelsResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}
	models := make([]string, 0, len(payload.Data))
	for _, model := range payload.Data {
		if model.ID == "" {
			continue
		}
		if strings.HasPrefix(model.ID, "gpt-") || strings.HasPrefix(model.ID, "o") {
			models = append(models, model.ID)
		}
	}
	return models, nil
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
	IsSeries      bool     `json:"isSeries"`
	SeriesName    string   `json:"seriesName"`
	VolumeNumber  *int     `json:"volumeNumber"`
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

func isValidEmail(value string) bool {
	value = strings.TrimSpace(value)
	if value == "" {
		return false
	}
	parts := strings.Split(value, "@")
	if len(parts) != 2 {
		return false
	}
	if parts[0] == "" || parts[1] == "" {
		return false
	}
	return strings.Contains(parts[1], ".")
}

func containsAuthor(authors []string, query string) bool {
	if query == "" {
		return false
	}
	for _, author := range authors {
		if strings.Contains(strings.ToLower(author), query) {
			return true
		}
	}
	return false
}
