package repository

import "book_manager/backend/internal/domain"

type BookRepository interface {
	Create(book domain.Book) error
	FindByID(id string) (domain.Book, bool)
	FindByISBN(isbn string) (domain.Book, bool)
	List() []domain.Book
}
