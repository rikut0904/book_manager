package userbooks

import (
	"errors"

	"book_manager/backend/internal/domain"
	"book_manager/backend/internal/idgen"
	"book_manager/backend/internal/repository"
)

var ErrUserBookExists = errors.New("user book already exists")

type Service struct {
	repo repository.UserBookRepository
}

type UpdateInput struct {
	Note         *string
	AcquiredAt   *string
	SeriesID     *string
	VolumeNumber *int
	SeriesSource *string
}

func NewService(repo repository.UserBookRepository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Create(userID, bookID, note, acquiredAt string) (domain.UserBook, error) {
	userBook := domain.UserBook{
		ID:         idgen.NewUserBook(),
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

func (s *Service) ListAll() []domain.UserBook {
	return s.repo.ListAll()
}

func (s *Service) ListBySeriesID(seriesID string) []domain.UserBook {
	return s.repo.ListBySeriesID(seriesID)
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
	if input.SeriesID != nil {
		userBook.SeriesID = *input.SeriesID
	}
	if input.VolumeNumber != nil {
		userBook.VolumeNumber = *input.VolumeNumber
	}
	if input.SeriesSource != nil {
		userBook.SeriesSource = *input.SeriesSource
	}
	if !s.repo.Update(userBook) {
		return domain.UserBook{}, false
	}
	return userBook, true
}

func (s *Service) Delete(id string) bool {
	return s.repo.Delete(id)
}
