package repository

import (
	"errors"
	"sync"

	"book_manager/backend/internal/domain"
)

var ErrBookTagExists = errors.New("book tag already exists")

type MemoryBookTagRepository struct {
	mu      sync.RWMutex
	items   map[string]domain.BookTag
	byUser  map[string][]string
	byTagID map[string][]string
}

func NewMemoryBookTagRepository() *MemoryBookTagRepository {
	return &MemoryBookTagRepository{
		items:   make(map[string]domain.BookTag),
		byUser:  make(map[string][]string),
		byTagID: make(map[string][]string),
	}
}

func (r *MemoryBookTagRepository) Create(item domain.BookTag) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	key := bookTagKey(item)
	if _, ok := r.items[key]; ok {
		return ErrBookTagExists
	}
	r.items[key] = item
	r.byUser[item.UserID] = append(r.byUser[item.UserID], key)
	r.byTagID[item.TagID] = append(r.byTagID[item.TagID], key)
	return nil
}

func (r *MemoryBookTagRepository) Delete(item domain.BookTag) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	key := bookTagKey(item)
	if _, ok := r.items[key]; !ok {
		return false
	}
	delete(r.items, key)
	if keys, ok := r.byUser[item.UserID]; ok {
		r.byUser[item.UserID] = removeID(keys, key)
	}
	if keys, ok := r.byTagID[item.TagID]; ok {
		r.byTagID[item.TagID] = removeID(keys, key)
	}
	return true
}

func (r *MemoryBookTagRepository) ListByUser(userID string) []domain.BookTag {
	r.mu.RLock()
	defer r.mu.RUnlock()

	keys := r.byUser[userID]
	items := make([]domain.BookTag, 0, len(keys))
	for _, key := range keys {
		if item, ok := r.items[key]; ok {
			items = append(items, item)
		}
	}
	return items
}

func (r *MemoryBookTagRepository) DeleteByTagID(tagID string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	keys := r.byTagID[tagID]
	for _, key := range keys {
		item, ok := r.items[key]
		if !ok {
			continue
		}
		delete(r.items, key)
		if userKeys, ok := r.byUser[item.UserID]; ok {
			r.byUser[item.UserID] = removeID(userKeys, key)
		}
	}
	delete(r.byTagID, tagID)
}

func bookTagKey(item domain.BookTag) string {
	return item.UserID + "::" + item.BookID + "::" + item.TagID
}
