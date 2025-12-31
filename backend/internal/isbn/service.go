package isbn

import (
	"errors"

	"book_manager/backend/internal/domain"
)

var ErrNotFound = errors.New("book not found")

type Service struct {
	books map[string]domain.Book
}

func NewService() *Service {
	return &Service{
		books: map[string]domain.Book{
			"9780000000000": {
				ID:            "book_demo_1",
				ISBN13:        "9780000000000",
				Title:         "サンプル書籍",
				Authors:       []string{"著者A"},
				Publisher:     "Book Manager Press",
				PublishedDate: "2024-01-01",
				ThumbnailURL:  "",
				Source:        "manual",
			},
		},
	}
}

func (s *Service) Lookup(isbn string) (domain.Book, error) {
	if book, ok := s.books[isbn]; ok {
		return book, nil
	}
	return domain.Book{}, ErrNotFound
}
