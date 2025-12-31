package repository

import (
	"errors"
	"sync"

	"book_manager/backend/internal/domain"
)

var ErrUserExists = errors.New("user already exists")

type MemoryUserRepository struct {
	mu      sync.RWMutex
	byID    map[string]domain.User
	byEmail map[string]domain.User
}

func NewMemoryUserRepository() *MemoryUserRepository {
	return &MemoryUserRepository{
		byID:    make(map[string]domain.User),
		byEmail: make(map[string]domain.User),
	}
}

func (r *MemoryUserRepository) Create(user domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.byEmail[user.Email]; ok {
		return ErrUserExists
	}
	r.byID[user.ID] = user
	r.byEmail[user.Email] = user
	return nil
}

func (r *MemoryUserRepository) FindByEmail(email string) (domain.User, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, ok := r.byEmail[email]
	return user, ok
}

func (r *MemoryUserRepository) FindByID(id string) (domain.User, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, ok := r.byID[id]
	return user, ok
}
