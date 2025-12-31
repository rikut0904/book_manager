package repository

import (
	"sync"

	"book_manager/backend/internal/domain"
)

type MemoryRecommendationRepository struct {
	mu      sync.RWMutex
	byID    map[string]domain.Recommendation
	ordered []string
}

func NewMemoryRecommendationRepository() *MemoryRecommendationRepository {
	return &MemoryRecommendationRepository{
		byID: make(map[string]domain.Recommendation),
	}
}

func (r *MemoryRecommendationRepository) Create(item domain.Recommendation) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.byID[item.ID] = item
	r.ordered = append(r.ordered, item.ID)
	return nil
}

func (r *MemoryRecommendationRepository) List() []domain.Recommendation {
	r.mu.RLock()
	defer r.mu.RUnlock()

	items := make([]domain.Recommendation, 0, len(r.ordered))
	for i := len(r.ordered) - 1; i >= 0; i-- {
		if item, ok := r.byID[r.ordered[i]]; ok {
			items = append(items, item)
		}
	}
	return items
}

func (r *MemoryRecommendationRepository) Delete(id string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.byID[id]; !ok {
		return false
	}
	delete(r.byID, id)
	r.ordered = removeID(r.ordered, id)
	return true
}
