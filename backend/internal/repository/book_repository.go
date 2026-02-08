package repository

import "book_manager/backend/internal/domain"

type BookRepository interface {
	Create(book domain.Book) error
	FindByID(id string) (domain.Book, bool)
	FindByISBN(isbn string) (domain.Book, bool)
	List() []domain.Book
	ListByUser(userID string) []domain.Book
	ListByIDs(ids []string) []domain.Book
	Delete(id string) bool
	Update(book domain.Book) bool
}
