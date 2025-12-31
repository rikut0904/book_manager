package repository

import "book_manager/backend/internal/domain"

type ProfileSettingsRepository interface {
	Get(userID string) (domain.ProfileSettings, bool)
	Upsert(settings domain.ProfileSettings)
	Delete(userID string)
}
