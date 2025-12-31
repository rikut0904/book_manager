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
	"book_manager/backend/internal/config"
	"book_manager/backend/internal/handler"
	"book_manager/backend/internal/isbn"
	"book_manager/backend/internal/repository"
	"book_manager/backend/internal/router"
)

func main() {
	cfg := config.Load()
	userRepo := repository.NewMemoryUserRepository()
	authService := auth.NewService(userRepo)
	isbnService := isbn.NewService(cfg.GoogleBooksBaseURL, cfg.GoogleBooksAPIKey)
	h := handler.New(authService, isbnService)
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
