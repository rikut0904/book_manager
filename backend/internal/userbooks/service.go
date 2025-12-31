package userbooks

import (
	"crypto/rand"
	"encoding/base64"
	"errors"

	"book_manager/backend/internal/domain"
	"book_manager/backend/internal/repository"
)

var ErrUserBookExists = errors.New("user book already exists")

type Service struct {
	repo repository.UserBookRepository
}

type UpdateInput struct {
	Note       *string
	AcquiredAt *string
}

func NewService(repo repository.UserBookRepository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Create(userID, bookID, note, acquiredAt string) (domain.UserBook, error) {
	userBook := domain.UserBook{
		ID:         newID(),
		UserID:     userID,
		BookID:     bookID,
		Note:       note,
		AcquiredAt: acquiredAt,
	}
	if err := s.repo.Create(userBook); err != nil {
		if errors.Is(err, repository.ErrUserBookExists) {
			return domain.UserBook{}, ErrUserBookExists
		}
		return domain.UserBook{}, err
	}
	return userBook, nil
}

func (s *Service) ListByUser(userID string) []domain.UserBook {
	return s.repo.ListByUser(userID)
}

func (s *Service) Update(id string, input UpdateInput) (domain.UserBook, bool) {
	userBook, ok := s.repo.FindByID(id)
	if !ok {
		return domain.UserBook{}, false
	}
	if input.Note != nil {
		userBook.Note = *input.Note
	}
	if input.AcquiredAt != nil {
		userBook.AcquiredAt = *input.AcquiredAt
	}
	if !s.repo.Update(userBook) {
		return domain.UserBook{}, false
	}
	return userBook, true
}

func (s *Service) Delete(id string) bool {
	return s.repo.Delete(id)
}

func newID() string {
	seed := make([]byte, 16)
	if _, err := rand.Read(seed); err != nil {
		return "user_book_demo"
	}
	return base64.RawURLEncoding.EncodeToString(seed)
}
