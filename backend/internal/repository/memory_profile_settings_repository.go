package repository

import (
	"sync"

	"book_manager/backend/internal/domain"
)

type MemoryProfileSettingsRepository struct {
	mu   sync.RWMutex
	byID map[string]domain.ProfileSettings
}

func NewMemoryProfileSettingsRepository() *MemoryProfileSettingsRepository {
	return &MemoryProfileSettingsRepository{
		byID: make(map[string]domain.ProfileSettings),
	}
}

func (r *MemoryProfileSettingsRepository) Get(userID string) (domain.ProfileSettings, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	settings, ok := r.byID[userID]
	return settings, ok
}

func (r *MemoryProfileSettingsRepository) Upsert(settings domain.ProfileSettings) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.byID[settings.UserID] = settings
}

func (r *MemoryProfileSettingsRepository) Delete(userID string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.byID, userID)
}
