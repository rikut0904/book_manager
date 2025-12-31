package gormrepo

import (
	"encoding/json"
	"time"

	"book_manager/backend/internal/domain"
	"book_manager/backend/internal/repository"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type IsbnCacheRepository struct {
	db *gorm.DB
}

func NewIsbnCacheRepository(db *gorm.DB) *IsbnCacheRepository {
	return &IsbnCacheRepository{db: db}
}

func (r *IsbnCacheRepository) Get(isbn string) (domain.IsbnCache, bool) {
	var model IsbnCache
	if err := r.db.First(&model, "isbn13 = ?", isbn).Error; err != nil {
		return domain.IsbnCache{}, false
	}
	var book domain.Book
	if err := json.Unmarshal(model.Payload, &book); err != nil {
		return domain.IsbnCache{}, false
	}
	return domain.IsbnCache{
		ISBN13:    model.ISBN13,
		Book:      book,
		FetchedAt: model.FetchedAt,
	}, true
}

func (r *IsbnCacheRepository) Upsert(entry domain.IsbnCache) error {
	payload, err := json.Marshal(entry.Book)
	if err != nil {
		return err
	}
	model := IsbnCache{
		ISBN13:    entry.ISBN13,
		Payload:   datatypes.JSON(payload),
		FetchedAt: entry.FetchedAt,
	}
	return r.db.Save(&model).Error
}

func (r *IsbnCacheRepository) DeleteBefore(t time.Time) error {
	return r.db.Delete(&IsbnCache{}, "fetched_at < ?", t).Error
}

var _ repository.IsbnCacheRepository = (*IsbnCacheRepository)(nil)
