package repository

import (
	"errors"
	"sync"

	"book_manager/backend/internal/domain"
)

var ErrFavoriteExists = errors.New("favorite already exists")

type MemoryFavoriteRepository struct {
	mu        sync.RWMutex
	byID      map[string]domain.Favorite
	byUser    map[string][]string
	byContent map[string]string
}

func NewMemoryFavoriteRepository() *MemoryFavoriteRepository {
	return &MemoryFavoriteRepository{
		byID:      make(map[string]domain.Favorite),
		byUser:    make(map[string][]string),
		byContent: make(map[string]string),
	}
}

func (r *MemoryFavoriteRepository) Create(favorite domain.Favorite) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	key := favoriteKey(favorite)
	if _, ok := r.byContent[key]; ok {
		return ErrFavoriteExists
	}
	r.byID[favorite.ID] = favorite
	r.byContent[key] = favorite.ID
	r.byUser[favorite.UserID] = append(r.byUser[favorite.UserID], favorite.ID)
	return nil
}

func (r *MemoryFavoriteRepository) ListByUser(userID string) []domain.Favorite {
	r.mu.RLock()
	defer r.mu.RUnlock()

	ids := r.byUser[userID]
	items := make([]domain.Favorite, 0, len(ids))
	for _, id := range ids {
		if fav, ok := r.byID[id]; ok {
			items = append(items, fav)
		}
	}
	return items
}

func (r *MemoryFavoriteRepository) ListBySeriesID(seriesID string) []domain.Favorite {
	r.mu.RLock()
	defer r.mu.RUnlock()

	items := make([]domain.Favorite, 0)
	for _, fav := range r.byID {
		if fav.SeriesID == seriesID {
			items = append(items, fav)
		}
	}
	return items
}

func (r *MemoryFavoriteRepository) Delete(id string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	fav, ok := r.byID[id]
	if !ok {
		return false
	}
	delete(r.byID, id)
	delete(r.byContent, favoriteKey(fav))
	if ids, ok := r.byUser[fav.UserID]; ok {
		r.byUser[fav.UserID] = removeID(ids, id)
	}
	return true
}

func favoriteKey(fav domain.Favorite) string {
	if fav.Type == "series" {
		return fav.UserID + ":series:" + fav.SeriesID
	}
	return fav.UserID + ":book:" + fav.BookID
}
