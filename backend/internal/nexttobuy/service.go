package nexttobuy

import (
	"book_manager/backend/internal/domain"
	"book_manager/backend/internal/idgen"
	"book_manager/backend/internal/repository"
)

type Service struct {
	repo repository.NextToBuyRepository
}

type UpdateInput struct {
	Title        *string
	SeriesName   *string
	VolumeNumber *int
	Note         *string
}

func NewService(repo repository.NextToBuyRepository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Create(userID, title, seriesName string, volumeNumber int, note string) (domain.NextToBuyManual, error) {
	item := domain.NextToBuyManual{
		ID:           idgen.NewNextToBuy(),
		UserID:       userID,
		Title:        title,
		SeriesName:   seriesName,
		VolumeNumber: volumeNumber,
		Note:         note,
	}
	if err := s.repo.Create(item); err != nil {
		return domain.NextToBuyManual{}, err
	}
	return item, nil
}

func (s *Service) ListByUser(userID string) []domain.NextToBuyManual {
	return s.repo.ListByUser(userID)
}

func (s *Service) Update(id string, input UpdateInput) (domain.NextToBuyManual, bool) {
	item, ok := s.repo.FindByID(id)
	if !ok {
		return domain.NextToBuyManual{}, false
	}
	if input.Title != nil {
		item.Title = *input.Title
	}
	if input.SeriesName != nil {
		item.SeriesName = *input.SeriesName
	}
	if input.VolumeNumber != nil {
		item.VolumeNumber = *input.VolumeNumber
	}
	if input.Note != nil {
		item.Note = *input.Note
	}
	if !s.repo.Update(item) {
		return domain.NextToBuyManual{}, false
	}
	return item, true
}

func (s *Service) Delete(id string) bool {
	return s.repo.Delete(id)
}
