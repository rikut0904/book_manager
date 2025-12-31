package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"book_manager/backend/internal/auth"
	"book_manager/backend/internal/books"
	"book_manager/backend/internal/config"
	"book_manager/backend/internal/favorites"
	"book_manager/backend/internal/follows"
	"book_manager/backend/internal/handler"
	"book_manager/backend/internal/isbn"
	"book_manager/backend/internal/nexttobuy"
	"book_manager/backend/internal/recommendations"
	"book_manager/backend/internal/reports"
	"book_manager/backend/internal/repository"
	"book_manager/backend/internal/router"
	"book_manager/backend/internal/tags"
	"book_manager/backend/internal/userbooks"
	"book_manager/backend/internal/users"
)

func main() {
	cfg := config.Load()
	userRepo := repository.NewMemoryUserRepository()
	bookRepo := repository.NewMemoryBookRepository()
	userBookRepo := repository.NewMemoryUserBookRepository()
	profileRepo := repository.NewMemoryProfileSettingsRepository()
	favoriteRepo := repository.NewMemoryFavoriteRepository()
	nextToBuyRepo := repository.NewMemoryNextToBuyRepository()
	tagRepo := repository.NewMemoryTagRepository()
	bookTagRepo := repository.NewMemoryBookTagRepository()
	recommendationRepo := repository.NewMemoryRecommendationRepository()
	authService := auth.NewService(userRepo)
	isbnService := isbn.NewService(cfg.GoogleBooksBaseURL, cfg.GoogleBooksAPIKey)
	bookService := books.NewService(bookRepo)
	userBookService := userbooks.NewService(userBookRepo)
	usersService := users.NewService(userRepo, profileRepo)
	followsService := follows.NewService()
	favoritesService := favorites.NewService(favoriteRepo)
	nextToBuyService := nexttobuy.NewService(nextToBuyRepo)
	tagsService := tags.NewService(tagRepo, bookTagRepo)
	recsService := recommendations.NewService(recommendationRepo)
	reportsService := reports.NewService(cfg.BookReportTo)
	_ = authService.SeedUser("user_demo", "demo@book.local", "demo", "password")
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
	)
	r := router.New(h)

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
