package repository

import (
	"time"

	"book_manager/backend/internal/domain"
)

type IsbnCacheRepository interface {
	Get(isbn string) (domain.IsbnCache, bool)
	Upsert(entry domain.IsbnCache) error
	DeleteBefore(t time.Time) error
}
