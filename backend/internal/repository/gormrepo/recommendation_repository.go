package gormrepo

import (
	"book_manager/backend/internal/domain"
	"book_manager/backend/internal/repository"
	"gorm.io/gorm"
)

type RecommendationRepository struct {
	db *gorm.DB
}

func NewRecommendationRepository(db *gorm.DB) *RecommendationRepository {
	return &RecommendationRepository{db: db}
}

func (r *RecommendationRepository) Create(item domain.Recommendation) error {
	model := Recommendation{
		ID:        item.ID,
		UserID:    item.UserID,
		BookID:    item.BookID,
		Comment:   item.Comment,
		CreatedAt: item.CreatedAt,
	}
	if err := r.db.Create(&model).Error; err != nil {
		return err
	}
	return nil
}

func (r *RecommendationRepository) List() []domain.Recommendation {
	var models []Recommendation
	if err := r.db.Order("created_at desc").Find(&models).Error; err != nil {
		return nil
	}
	items := make([]domain.Recommendation, 0, len(models))
	for _, model := range models {
		items = append(items, domain.Recommendation{
			ID:        model.ID,
			UserID:    model.UserID,
			BookID:    model.BookID,
			Comment:   model.Comment,
			CreatedAt: model.CreatedAt,
		})
	}
	return items
}

func (r *RecommendationRepository) Delete(id string) bool {
	if err := r.db.Delete(&Recommendation{}, "id = ?", id).Error; err != nil {
		return false
	}
	return true
}

var _ repository.RecommendationRepository = (*RecommendationRepository)(nil)
