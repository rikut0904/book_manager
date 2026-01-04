package gormrepo

import (
	"book_manager/backend/internal/domain"
	"book_manager/backend/internal/repository"
	"gorm.io/gorm"
)

type ProfileSettingsRepository struct {
	db *gorm.DB
}

func NewProfileSettingsRepository(db *gorm.DB) *ProfileSettingsRepository {
	return &ProfileSettingsRepository{db: db}
}

func (r *ProfileSettingsRepository) Get(userID string) (domain.ProfileSettings, bool) {
	var model ProfileSettings
	if err := r.db.First(&model, "user_id = ?", userID).Error; err != nil {
		return domain.ProfileSettings{}, false
	}
	return domain.ProfileSettings{
		UserID:        model.UserID,
		Visibility:    model.Visibility,
		OpenAIEnabled: model.OpenAIEnabled,
		OpenAIModel:   model.OpenAIModel,
		OpenAIAPIKey:  model.OpenAIAPIKey,
	}, true
}

func (r *ProfileSettingsRepository) Upsert(settings domain.ProfileSettings) {
	model := ProfileSettings{
		UserID:        settings.UserID,
		Visibility:    settings.Visibility,
		OpenAIEnabled: settings.OpenAIEnabled,
		OpenAIModel:   settings.OpenAIModel,
		OpenAIAPIKey:  settings.OpenAIAPIKey,
	}
	r.db.Save(&model)
}

func (r *ProfileSettingsRepository) Delete(userID string) {
	r.db.Delete(&ProfileSettings{}, "user_id = ?", userID)
}

var _ repository.ProfileSettingsRepository = (*ProfileSettingsRepository)(nil)
