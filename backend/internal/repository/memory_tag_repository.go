package repository

import (
	"errors"
	"sync"

	"book_manager/backend/internal/domain"
)

var ErrTagExists = errors.New("tag already exists")

type MemoryTagRepository struct {
	mu     sync.RWMutex
	byID   map[string]domain.Tag
	byUser map[string][]string
	byName map[string]string
}

func NewMemoryTagRepository() *MemoryTagRepository {
	return &MemoryTagRepository{
		byID:   make(map[string]domain.Tag),
		byUser: make(map[string][]string),
		byName: make(map[string]string),
	}
}

func (r *MemoryTagRepository) Create(tag domain.Tag) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	key := tag.OwnerUserID + "::" + tag.Name
	if _, ok := r.byName[key]; ok {
		return ErrTagExists
	}
	r.byID[tag.ID] = tag
	r.byName[key] = tag.ID
	r.byUser[tag.OwnerUserID] = append(r.byUser[tag.OwnerUserID], tag.ID)
	return nil
}

func (r *MemoryTagRepository) ListByUser(userID string) []domain.Tag {
	r.mu.RLock()
	defer r.mu.RUnlock()

	ids := r.byUser[userID]
	items := make([]domain.Tag, 0, len(ids))
	for _, id := range ids {
		if tag, ok := r.byID[id]; ok {
			items = append(items, tag)
		}
	}
	return items
}

func (r *MemoryTagRepository) FindByID(id string) (domain.Tag, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tag, ok := r.byID[id]
	return tag, ok
}

func (r *MemoryTagRepository) Delete(id string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	tag, ok := r.byID[id]
	if !ok {
		return false
	}
	delete(r.byID, id)
	delete(r.byName, tag.OwnerUserID+"::"+tag.Name)
	if ids, ok := r.byUser[tag.OwnerUserID]; ok {
		r.byUser[tag.OwnerUserID] = removeID(ids, id)
	}
	return true
}
