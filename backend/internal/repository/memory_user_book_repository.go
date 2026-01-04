package repository

import (
	"errors"
	"sync"

	"book_manager/backend/internal/domain"
)

var ErrUserBookExists = errors.New("user book already exists")

type MemoryUserBookRepository struct {
	mu         sync.RWMutex
	byID       map[string]domain.UserBook
	byUser     map[string][]string
	byUserBook map[string]string
}

func NewMemoryUserBookRepository() *MemoryUserBookRepository {
	return &MemoryUserBookRepository{
		byID:       make(map[string]domain.UserBook),
		byUser:     make(map[string][]string),
		byUserBook: make(map[string]string),
	}
}

func (r *MemoryUserBookRepository) Create(userBook domain.UserBook) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	key := userBookKey(userBook.UserID, userBook.BookID)
	if _, ok := r.byUserBook[key]; ok {
		return ErrUserBookExists
	}
	r.byID[userBook.ID] = userBook
	r.byUserBook[key] = userBook.ID
	r.byUser[userBook.UserID] = append(r.byUser[userBook.UserID], userBook.ID)
	return nil
}

func (r *MemoryUserBookRepository) ListByUser(userID string) []domain.UserBook {
	r.mu.RLock()
	defer r.mu.RUnlock()

	ids := r.byUser[userID]
	books := make([]domain.UserBook, 0, len(ids))
	for _, id := range ids {
		if book, ok := r.byID[id]; ok {
			books = append(books, book)
		}
	}
	return books
}

func (r *MemoryUserBookRepository) ListAll() []domain.UserBook {
	r.mu.RLock()
	defer r.mu.RUnlock()

	books := make([]domain.UserBook, 0, len(r.byID))
	for _, book := range r.byID {
		books = append(books, book)
	}
	return books
}

func (r *MemoryUserBookRepository) ListBySeriesID(seriesID string) []domain.UserBook {
	r.mu.RLock()
	defer r.mu.RUnlock()

	books := make([]domain.UserBook, 0)
	for _, book := range r.byID {
		if book.SeriesID == seriesID {
			books = append(books, book)
		}
	}
	return books
}

func (r *MemoryUserBookRepository) FindByID(id string) (domain.UserBook, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	book, ok := r.byID[id]
	return book, ok
}

func (r *MemoryUserBookRepository) Update(userBook domain.UserBook) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.byID[userBook.ID]; !ok {
		return false
	}
	r.byID[userBook.ID] = userBook
	return true
}

func (r *MemoryUserBookRepository) Delete(id string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	book, ok := r.byID[id]
	if !ok {
		return false
	}
	delete(r.byID, id)
	key := userBookKey(book.UserID, book.BookID)
	delete(r.byUserBook, key)
	if ids, ok := r.byUser[book.UserID]; ok {
		r.byUser[book.UserID] = removeID(ids, id)
	}
	return true
}

func userBookKey(userID, bookID string) string {
	return userID + "::" + bookID
}

func removeID(ids []string, id string) []string {
	result := ids[:0]
	for _, item := range ids {
		if item != id {
			result = append(result, item)
		}
	}
	return result
}
