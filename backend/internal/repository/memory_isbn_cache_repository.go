package repository

import (
	"sync"
	"time"

	"book_manager/backend/internal/domain"
)

type MemoryIsbnCacheRepository struct {
	mu    sync.RWMutex
	cache map[string]domain.IsbnCache
}

func NewMemoryIsbnCacheRepository() *MemoryIsbnCacheRepository {
	return &MemoryIsbnCacheRepository{
		cache: make(map[string]domain.IsbnCache),
	}
}

func (r *MemoryIsbnCacheRepository) Get(isbn string) (domain.IsbnCache, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	entry, ok := r.cache[isbn]
	return entry, ok
}

func (r *MemoryIsbnCacheRepository) Upsert(entry domain.IsbnCache) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.cache[entry.ISBN13] = entry
	return nil
}

func (r *MemoryIsbnCacheRepository) DeleteBefore(t time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for key, entry := range r.cache {
		if entry.FetchedAt.Before(t) {
			delete(r.cache, key)
		}
	}
	return nil
}
