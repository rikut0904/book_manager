package repository

import (
	"errors"
	"sync"

	"book_manager/backend/internal/domain"
)

var ErrBookExists = errors.New("book already exists")

type MemoryBookRepository struct {
	mu      sync.RWMutex
	byID    map[string]domain.Book
	byISBN  map[string]domain.Book
	ordered []string
}

func NewMemoryBookRepository() *MemoryBookRepository {
	return &MemoryBookRepository{
		byID:   make(map[string]domain.Book),
		byISBN: make(map[string]domain.Book),
	}
}

func (r *MemoryBookRepository) Create(book domain.Book) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if book.ISBN13 != "" {
		if _, ok := r.byISBN[book.ISBN13]; ok {
			return ErrBookExists
		}
	}
	r.byID[book.ID] = book
	if book.ISBN13 != "" {
		r.byISBN[book.ISBN13] = book
	}
	r.ordered = append(r.ordered, book.ID)
	return nil
}

func (r *MemoryBookRepository) FindByID(id string) (domain.Book, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	book, ok := r.byID[id]
	return book, ok
}

func (r *MemoryBookRepository) FindByISBN(isbn string) (domain.Book, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	book, ok := r.byISBN[isbn]
	return book, ok
}

func (r *MemoryBookRepository) List() []domain.Book {
	r.mu.RLock()
	defer r.mu.RUnlock()

	books := make([]domain.Book, 0, len(r.ordered))
	for _, id := range r.ordered {
		if book, ok := r.byID[id]; ok {
			books = append(books, book)
		}
	}
	return books
}

func (r *MemoryBookRepository) Delete(id string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	book, ok := r.byID[id]
	if !ok {
		return false
	}
	delete(r.byID, id)
	if book.ISBN13 != "" {
		delete(r.byISBN, book.ISBN13)
	}
	for i, storedID := range r.ordered {
		if storedID == id {
			r.ordered = append(r.ordered[:i], r.ordered[i+1:]...)
			break
		}
	}
	return true
}
