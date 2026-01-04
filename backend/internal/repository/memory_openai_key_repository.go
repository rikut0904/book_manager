package repository

import (
	"sync"

	"book_manager/backend/internal/domain"
)

type MemoryOpenAIKeyRepository struct {
	mu      sync.RWMutex
	byID    map[string]domain.OpenAIKey
	ordered []string
}

func NewMemoryOpenAIKeyRepository() *MemoryOpenAIKeyRepository {
	return &MemoryOpenAIKeyRepository{
		byID: make(map[string]domain.OpenAIKey),
	}
}

func (r *MemoryOpenAIKeyRepository) Create(key domain.OpenAIKey) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.byID[key.ID] = key
	r.ordered = append(r.ordered, key.ID)
	return nil
}

func (r *MemoryOpenAIKeyRepository) List() []domain.OpenAIKey {
	r.mu.RLock()
	defer r.mu.RUnlock()

	items := make([]domain.OpenAIKey, 0, len(r.ordered))
	for _, id := range r.ordered {
		if key, ok := r.byID[id]; ok {
			items = append(items, key)
		}
	}
	return items
}

func (r *MemoryOpenAIKeyRepository) Delete(id string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.byID[id]; !ok {
		return false
	}
	delete(r.byID, id)
	for i, storedID := range r.ordered {
		if storedID == id {
			r.ordered = append(r.ordered[:i], r.ordered[i+1:]...)
			break
		}
	}
	return true
}
