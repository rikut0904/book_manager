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
	ordered []string
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
	r.ordered = append(r.ordered, user.ID)
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

func (r *MemoryUserRepository) List() []domain.User {
	r.mu.RLock()
	defer r.mu.RUnlock()

	users := make([]domain.User, 0, len(r.ordered))
	for _, id := range r.ordered {
		if user, ok := r.byID[id]; ok {
			users = append(users, user)
		}
	}
	return users
}

func (r *MemoryUserRepository) Update(user domain.User) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	existing, ok := r.byID[user.ID]
	if !ok {
		return false
	}
	if existing.Email != user.Email {
		if _, ok := r.byEmail[user.Email]; ok {
			return false
		}
		delete(r.byEmail, existing.Email)
		r.byEmail[user.Email] = user
	} else {
		r.byEmail[user.Email] = user
	}
	r.byID[user.ID] = user
	return true
}

func (r *MemoryUserRepository) Delete(id string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	user, ok := r.byID[id]
	if !ok {
		return false
	}
	delete(r.byID, id)
	delete(r.byEmail, user.Email)
	r.ordered = removeUserID(r.ordered, id)
	return true
}

func removeUserID(ids []string, id string) []string {
	result := ids[:0]
	for _, item := range ids {
		if item != id {
			result = append(result, item)
		}
	}
	return result
}
