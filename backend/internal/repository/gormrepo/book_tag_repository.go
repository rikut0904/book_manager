package gormrepo

import (
	"book_manager/backend/internal/domain"
	"book_manager/backend/internal/repository"
	"gorm.io/gorm"
)

type BookTagRepository struct {
	db *gorm.DB
}

func NewBookTagRepository(db *gorm.DB) *BookTagRepository {
	return &BookTagRepository{db: db}
}

func (r *BookTagRepository) Create(item domain.BookTag) error {
	model := BookTag{
		UserID: item.UserID,
		BookID: item.BookID,
		TagID:  item.TagID,
	}
	if err := r.db.Create(&model).Error; err != nil {
		if isUniqueViolation(err) {
			return repository.ErrBookTagExists
		}
		return err
	}
	return nil
}

func (r *BookTagRepository) Delete(item domain.BookTag) bool {
	if err := r.db.Delete(&BookTag{}, "user_id = ? AND book_id = ? AND tag_id = ?", item.UserID, item.BookID, item.TagID).Error; err != nil {
		return false
	}
	return true
}

func (r *BookTagRepository) ListByUser(userID string) []domain.BookTag {
	var models []BookTag
	if err := r.db.Where("user_id = ?", userID).Find(&models).Error; err != nil {
		return nil
	}
	items := make([]domain.BookTag, 0, len(models))
	for _, model := range models {
		items = append(items, domain.BookTag{
			UserID: model.UserID,
			BookID: model.BookID,
			TagID:  model.TagID,
		})
	}
	return items
}

func (r *BookTagRepository) DeleteByTagID(tagID string) {
	r.db.Delete(&BookTag{}, "tag_id = ?", tagID)
}

var _ repository.BookTagRepository = (*BookTagRepository)(nil)
