package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"book_manager/backend/internal/auth"
	"book_manager/backend/internal/books"
	"book_manager/backend/internal/config"
	"book_manager/backend/internal/db"
	"book_manager/backend/internal/favorites"
	"book_manager/backend/internal/follows"
	"book_manager/backend/internal/openaikeys"
	"book_manager/backend/internal/handler"
	"book_manager/backend/internal/isbn"
	"book_manager/backend/internal/nexttobuy"
	"book_manager/backend/internal/recommendations"
	"book_manager/backend/internal/reports"
	"book_manager/backend/internal/repository"
	"book_manager/backend/internal/repository/gormrepo"
	"book_manager/backend/internal/router"
	"book_manager/backend/internal/series"
	"book_manager/backend/internal/tags"
	"book_manager/backend/internal/userbooks"
	"book_manager/backend/internal/users"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	_ = godotenv.Load("../.env")

	cfg := config.Load()
	aiPrompt := loadPrompt("prompt.md")
	var (
		userRepo           repository.UserRepository
		bookRepo           repository.BookRepository
		userBookRepo       repository.UserBookRepository
		profileRepo        repository.ProfileSettingsRepository
		favoriteRepo       repository.FavoriteRepository
		nextToBuyRepo      repository.NextToBuyRepository
		tagRepo            repository.TagRepository
		bookTagRepo        repository.BookTagRepository
		recommendationRepo repository.RecommendationRepository
		isbnCacheRepo      repository.IsbnCacheRepository
		auditLogRepo       repository.AuditLogRepository
		seriesRepo         repository.SeriesRepository
		openAIKeyRepo      repository.OpenAIKeyRepository
	)

	if cfg.DatabaseURL != "" {
		dbConn, err := db.Open(cfg.DatabaseURL)
		if err != nil {
			log.Fatalf("db connection error: %v", err)
		}
		if err := dbConn.AutoMigrate(
			&gormrepo.User{},
			&gormrepo.ProfileSettings{},
			&gormrepo.Book{},
			&gormrepo.UserBook{},
			&gormrepo.Favorite{},
			&gormrepo.NextToBuyManual{},
			&gormrepo.Tag{},
			&gormrepo.BookTag{},
			&gormrepo.Recommendation{},
			&gormrepo.IsbnCache{},
			&gormrepo.AuditLog{},
			&gormrepo.Series{},
			&gormrepo.OpenAIKey{},
		); err != nil {
			log.Fatalf("db migrate error: %v", err)
		}
		userRepo = gormrepo.NewUserRepository(dbConn)
		bookRepo = gormrepo.NewBookRepository(dbConn)
		userBookRepo = gormrepo.NewUserBookRepository(dbConn)
		profileRepo = gormrepo.NewProfileSettingsRepository(dbConn)
		favoriteRepo = gormrepo.NewFavoriteRepository(dbConn)
		nextToBuyRepo = gormrepo.NewNextToBuyRepository(dbConn)
		tagRepo = gormrepo.NewTagRepository(dbConn)
		bookTagRepo = gormrepo.NewBookTagRepository(dbConn)
		recommendationRepo = gormrepo.NewRecommendationRepository(dbConn)
		isbnCacheRepo = gormrepo.NewIsbnCacheRepository(dbConn)
		auditLogRepo = gormrepo.NewAuditLogRepository(dbConn)
		seriesRepo = gormrepo.NewSeriesRepository(dbConn)
		openAIKeyRepo = gormrepo.NewOpenAIKeyRepository(dbConn)
	} else {
		userRepo = repository.NewMemoryUserRepository()
		bookRepo = repository.NewMemoryBookRepository()
		userBookRepo = repository.NewMemoryUserBookRepository()
		profileRepo = repository.NewMemoryProfileSettingsRepository()
		favoriteRepo = repository.NewMemoryFavoriteRepository()
		nextToBuyRepo = repository.NewMemoryNextToBuyRepository()
		tagRepo = repository.NewMemoryTagRepository()
		bookTagRepo = repository.NewMemoryBookTagRepository()
		recommendationRepo = repository.NewMemoryRecommendationRepository()
		isbnCacheRepo = repository.NewMemoryIsbnCacheRepository()
		auditLogRepo = repository.NewMemoryAuditLogRepository()
		seriesRepo = repository.NewMemorySeriesRepository()
		openAIKeyRepo = repository.NewMemoryOpenAIKeyRepository()
	}
	authService := auth.NewService(userRepo)
	isbnCacheTTL := time.Duration(cfg.IsbnCacheTTLMinutes) * time.Minute
	isbnService := isbn.NewService(cfg.GoogleBooksBaseURL, cfg.GoogleBooksAPIKey, isbnCacheTTL, isbnCacheRepo)
	bookService := books.NewService(bookRepo)
	userBookService := userbooks.NewService(userBookRepo)
	usersService := users.NewService(userRepo, profileRepo)
	followsService := follows.NewService()
	favoritesService := favorites.NewService(favoriteRepo)
	nextToBuyService := nexttobuy.NewService(nextToBuyRepo)
	tagsService := tags.NewService(tagRepo, bookTagRepo)
	recsService := recommendations.NewService(recommendationRepo)
	reportsService := reports.NewService(cfg.BookReportTo, reports.SMTPConfig{
		Host: cfg.SMTPHost,
		Port: cfg.SMTPPort,
		User: cfg.SMTPUser,
		Pass: cfg.SMTPPass,
		From: cfg.SMTPFrom,
	})
	seriesService := series.NewService(seriesRepo)
	openAIKeyService := openaikeys.NewService(openAIKeyRepo)
	_ = authService.SeedUser("user_demo", "demo@book.local", "demo", "password")
	if count := normalizeBooks(bookService); count > 0 {
		log.Printf("normalized %d book titles", count)
	}
	if count := seriesService.NormalizeAll(); count > 0 {
		log.Printf("normalized %d series names", count)
	}
	h := handler.New(
		authService,
		isbnService,
		bookService,
		userBookService,
		usersService,
		followsService,
		favoritesService,
		nextToBuyService,
		tagsService,
		recsService,
		reportsService,
		seriesService,
		openAIKeyService,
		cfg.OpenAIAPIKey,
		cfg.OpenAIDefaultModel,
		aiPrompt,
		parseAdminUserIDs(cfg.AdminUserIDs),
	)
	r := router.New(h, auditLogRepo, cfg.CORSAllowedOrigins)

	server := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Printf("api server starting on :%s (env=%s)", cfg.Port, cfg.Env)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	go startAuditCleanup(auditLogRepo)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Println("shutting down server...")
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("shutdown error: %v", err)
	}
}

func loadPrompt(path string) string {
	data, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		log.Printf("prompt load error: %v", err)
		return ""
	}
	return string(data)
}

func parseAdminUserIDs(value string) []string {
	if strings.TrimSpace(value) == "" {
		return nil
	}
	parts := strings.Split(value, ",")
	ids := make([]string, 0, len(parts))
	for _, part := range parts {
		id := strings.TrimSpace(part)
		if id == "" {
			continue
		}
		ids = append(ids, id)
	}
	return ids
}

func normalizeBooks(bookService *books.Service) int {
	items := bookService.List()
	updated := 0
	for _, item := range items {
		cleaned := isbn.NormalizeTitle(item.Title)
		if cleaned == "" || cleaned == item.Title {
			continue
		}
		item.Title = cleaned
		if bookService.Update(item) {
			updated++
		}
	}
	return updated
}


func startAuditCleanup(repo repository.AuditLogRepository) {
	if repo == nil {
		return
	}
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	for {
		cutoff := time.Now().AddDate(0, 0, -90)
		_ = repo.DeleteBefore(cutoff)
		<-ticker.C
	}
}
