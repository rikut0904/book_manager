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
	SMTPHost            string
	SMTPPort            string
	SMTPUser            string
	SMTPPass            string
	SMTPFrom            string
	CORSAllowedOrigins  string
	OpenAIAPIKey        string
	OpenAIDefaultModel  string
	AdminUserIDs        string
	FirebaseProjectID   string
	FirebaseAPIKey      string
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
		SMTPHost:            getEnv("SMTP_HOST", ""),
		SMTPPort:            getEnv("SMTP_PORT", "587"),
		SMTPUser:            getEnv("SMTP_USER", ""),
		SMTPPass:            getEnv("SMTP_PASS", ""),
		SMTPFrom:            getEnv("SMTP_FROM", ""),
		CORSAllowedOrigins:  getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:3000"),
		OpenAIAPIKey:        getEnv("OPENAI_API_KEY", ""),
		OpenAIDefaultModel:  getEnv("OPENAI_DEFAULT_MODEL", "gpt-4o-mini"),
		AdminUserIDs:        getEnv("ADMIN_USER_IDS", ""),
		FirebaseProjectID:   getEnv("FIREBASE_PROJECT_ID", ""),
		FirebaseAPIKey:      getEnv("FIREBASE_API_KEY", ""),
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
