package repository

import (
	"sync"

	"book_manager/backend/internal/domain"
)

type MemoryAdminUserRepository struct {
	mu       sync.RWMutex
	byUserID map[string]domain.AdminUser
	ordered  []string
}

func NewMemoryAdminUserRepository() *MemoryAdminUserRepository {
	return &MemoryAdminUserRepository{
		byUserID: make(map[string]domain.AdminUser),
	}
}

func (r *MemoryAdminUserRepository) Create(adminUser domain.AdminUser) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.byUserID[adminUser.UserID] = adminUser
	r.ordered = append(r.ordered, adminUser.UserID)
	return nil
}

func (r *MemoryAdminUserRepository) FindByUserID(userID string) (domain.AdminUser, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	adminUser, ok := r.byUserID[userID]
	return adminUser, ok
}

func (r *MemoryAdminUserRepository) List() []domain.AdminUser {
	r.mu.RLock()
	defer r.mu.RUnlock()

	items := make([]domain.AdminUser, 0, len(r.ordered))
	for _, userID := range r.ordered {
		if adminUser, ok := r.byUserID[userID]; ok {
			items = append(items, adminUser)
		}
	}
	return items
}

func (r *MemoryAdminUserRepository) Delete(userID string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.byUserID[userID]; !ok {
		return false
	}
	delete(r.byUserID, userID)
	for i, storedUserID := range r.ordered {
		if storedUserID == userID {
			r.ordered = append(r.ordered[:i], r.ordered[i+1:]...)
			break
		}
	}
	return true
}
