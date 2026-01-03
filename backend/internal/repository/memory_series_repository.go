package repository

import (
	"sync"

	"book_manager/backend/internal/domain"
)

type MemorySeriesRepository struct {
	mu      sync.RWMutex
	byID    map[string]domain.Series
	byName  map[string]string
	ordered []string
}

func NewMemorySeriesRepository() *MemorySeriesRepository {
	return &MemorySeriesRepository{
		byID:   make(map[string]domain.Series),
		byName: make(map[string]string),
	}
}

func (r *MemorySeriesRepository) Create(series domain.Series) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.byName[series.NormalizedName]; ok {
		return ErrSeriesExists
	}
	r.byID[series.ID] = series
	r.byName[series.NormalizedName] = series.ID
	r.ordered = append(r.ordered, series.ID)
	return nil
}

func (r *MemorySeriesRepository) List() []domain.Series {
	r.mu.RLock()
	defer r.mu.RUnlock()

	items := make([]domain.Series, 0, len(r.ordered))
	for _, id := range r.ordered {
		if item, ok := r.byID[id]; ok {
			items = append(items, item)
		}
	}
	return items
}

func (r *MemorySeriesRepository) FindByNormalizedName(name string) (domain.Series, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if id, ok := r.byName[name]; ok {
		item, exists := r.byID[id]
		return item, exists
	}
	return domain.Series{}, false
}

func (r *MemorySeriesRepository) Delete(id string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	item, ok := r.byID[id]
	if !ok {
		return false
	}
	delete(r.byID, id)
	if item.NormalizedName != "" {
		delete(r.byName, item.NormalizedName)
	}
	for i, storedID := range r.ordered {
		if storedID == id {
			r.ordered = append(r.ordered[:i], r.ordered[i+1:]...)
			break
		}
	}
	return true
}

func (r *MemorySeriesRepository) Update(series domain.Series) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	current, ok := r.byID[series.ID]
	if !ok {
		return false
	}
	if current.NormalizedName != "" && current.NormalizedName != series.NormalizedName {
		delete(r.byName, current.NormalizedName)
	}
	r.byID[series.ID] = series
	if series.NormalizedName != "" {
		r.byName[series.NormalizedName] = series.ID
	}
	return true
}
