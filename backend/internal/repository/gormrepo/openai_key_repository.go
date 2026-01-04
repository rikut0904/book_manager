package gormrepo

import (
	"book_manager/backend/internal/domain"
	"book_manager/backend/internal/repository"
	"gorm.io/gorm"
)

type OpenAIKeyRepository struct {
	db *gorm.DB
}

func NewOpenAIKeyRepository(db *gorm.DB) *OpenAIKeyRepository {
	return &OpenAIKeyRepository{db: db}
}

func (r *OpenAIKeyRepository) Create(key domain.OpenAIKey) error {
	model := OpenAIKey{
		ID:        key.ID,
		Name:      key.Name,
		APIKey:    key.APIKey,
		CreatedAt: key.CreatedAt,
	}
	return r.db.Create(&model).Error
}

func (r *OpenAIKeyRepository) List() []domain.OpenAIKey {
	var models []OpenAIKey
	if err := r.db.Order("created_at asc").Find(&models).Error; err != nil {
		return nil
	}
	items := make([]domain.OpenAIKey, 0, len(models))
	for _, model := range models {
		items = append(items, domain.OpenAIKey{
			ID:        model.ID,
			Name:      model.Name,
			APIKey:    model.APIKey,
			CreatedAt: model.CreatedAt,
		})
	}
	return items
}

func (r *OpenAIKeyRepository) Delete(id string) bool {
	if err := r.db.Delete(&OpenAIKey{}, "id = ?", id).Error; err != nil {
		return false
	}
	return true
}

var _ repository.OpenAIKeyRepository = (*OpenAIKeyRepository)(nil)
