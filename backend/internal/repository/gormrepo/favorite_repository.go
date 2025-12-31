package gormrepo

import (
	"book_manager/backend/internal/domain"
	"book_manager/backend/internal/repository"
	"gorm.io/gorm"
)

type FavoriteRepository struct {
	db *gorm.DB
}

func NewFavoriteRepository(db *gorm.DB) *FavoriteRepository {
	return &FavoriteRepository{db: db}
}

func (r *FavoriteRepository) Create(favorite domain.Favorite) error {
	model := Favorite{
		ID:     favorite.ID,
		UserID: favorite.UserID,
		Type:   favorite.Type,
	}
	if favorite.Type == "series" {
		model.SeriesID = stringPtr(favorite.SeriesID)
	} else {
		model.BookID = stringPtr(favorite.BookID)
	}
	if err := r.db.Create(&model).Error; err != nil {
		if isUniqueViolation(err) {
			return repository.ErrFavoriteExists
		}
		return err
	}
	return nil
}

func (r *FavoriteRepository) ListByUser(userID string) []domain.Favorite {
	var models []Favorite
	if err := r.db.Where("user_id = ?", userID).Order("id asc").Find(&models).Error; err != nil {
		return nil
	}
	items := make([]domain.Favorite, 0, len(models))
	for _, model := range models {
		items = append(items, domain.Favorite{
			ID:       model.ID,
			UserID:   model.UserID,
			Type:     model.Type,
			BookID:   valueOrEmptyString(model.BookID),
			SeriesID: valueOrEmptyString(model.SeriesID),
		})
	}
	return items
}

func (r *FavoriteRepository) Delete(id string) bool {
	if err := r.db.Delete(&Favorite{}, "id = ?", id).Error; err != nil {
		return false
	}
	return true
}

func stringPtr(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}

func valueOrEmptyString(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

var _ repository.FavoriteRepository = (*FavoriteRepository)(nil)
