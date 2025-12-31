package repository

import (
	"sync"

	"book_manager/backend/internal/domain"
)

type MemoryNextToBuyRepository struct {
	mu     sync.RWMutex
	byID   map[string]domain.NextToBuyManual
	byUser map[string][]string
}

func NewMemoryNextToBuyRepository() *MemoryNextToBuyRepository {
	return &MemoryNextToBuyRepository{
		byID:   make(map[string]domain.NextToBuyManual),
		byUser: make(map[string][]string),
	}
}

func (r *MemoryNextToBuyRepository) Create(item domain.NextToBuyManual) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.byID[item.ID] = item
	r.byUser[item.UserID] = append(r.byUser[item.UserID], item.ID)
	return nil
}

func (r *MemoryNextToBuyRepository) ListByUser(userID string) []domain.NextToBuyManual {
	r.mu.RLock()
	defer r.mu.RUnlock()

	ids := r.byUser[userID]
	items := make([]domain.NextToBuyManual, 0, len(ids))
	for _, id := range ids {
		if item, ok := r.byID[id]; ok {
			items = append(items, item)
		}
	}
	return items
}

func (r *MemoryNextToBuyRepository) FindByID(id string) (domain.NextToBuyManual, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	item, ok := r.byID[id]
	return item, ok
}

func (r *MemoryNextToBuyRepository) Update(item domain.NextToBuyManual) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.byID[item.ID]; !ok {
		return false
	}
	r.byID[item.ID] = item
	return true
}

func (r *MemoryNextToBuyRepository) Delete(id string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	item, ok := r.byID[id]
	if !ok {
		return false
	}
	delete(r.byID, id)
	if ids, ok := r.byUser[item.UserID]; ok {
		r.byUser[item.UserID] = removeID(ids, id)
	}
	return true
}
