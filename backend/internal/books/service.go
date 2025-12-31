package books

import (
	"crypto/rand"
	"encoding/base64"
	"errors"

	"book_manager/backend/internal/domain"
	"book_manager/backend/internal/repository"
)

var ErrBookExists = errors.New("book already exists")

type Service struct {
	books repository.BookRepository
}

func NewService(books repository.BookRepository) *Service {
	return &Service{
		books: books,
	}
}

func (s *Service) Create(input domain.Book) (domain.Book, error) {
	book := input
	book.ID = newID()
	if book.Source == "" {
		book.Source = "manual"
	}
	if err := s.books.Create(book); err != nil {
		if errors.Is(err, repository.ErrBookExists) {
			return domain.Book{}, ErrBookExists
		}
		return domain.Book{}, err
	}
	return book, nil
}

func (s *Service) List() []domain.Book {
	return s.books.List()
}

func (s *Service) Get(id string) (domain.Book, bool) {
	return s.books.FindByID(id)
}

func (s *Service) FindByISBN(isbn string) (domain.Book, bool) {
	return s.books.FindByISBN(isbn)
}

func newID() string {
	seed := make([]byte, 16)
	if _, err := rand.Read(seed); err != nil {
		return "book_demo"
	}
	return base64.RawURLEncoding.EncodeToString(seed)
}
