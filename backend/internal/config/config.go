package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port                string
	Env                 string
	GoogleBooksAPIKey   string
	GoogleBooksBaseURL  string
	BookReportTo        string
	IsbnCacheTTLMinutes int
	DatabaseURL         string
}

func Load() Config {
	return Config{
		Port:                getEnv("PORT", "8080"),
		Env:                 getEnv("APP_ENV", "local"),
		GoogleBooksAPIKey:   getEnv("GOOGLE_BOOKS_API_KEY", ""),
		GoogleBooksBaseURL:  getEnv("GOOGLE_BOOKS_BASE_URL", "https://www.googleapis.com/books/v1/volumes"),
		BookReportTo:        getEnv("BOOK_REPORT_TO", "product@rikut0904.site"),
		IsbnCacheTTLMinutes: getEnvInt("ISBN_CACHE_TTL_MINUTES", 1440),
		DatabaseURL:         getEnv("DATABASE_URL", ""),
	}
}

func getEnvInt(key string, fallback int) int {
	if value, ok := os.LookupEnv(key); ok {
		if parsed, err := strconv.Atoi(value); err == nil {
			return parsed
		}
	}
	return fallback
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
