package favorites

import (
	"errors"

	"book_manager/backend/internal/domain"
	"book_manager/backend/internal/idgen"
	"book_manager/backend/internal/repository"
)

var ErrFavoriteExists = errors.New("favorite already exists")

type Service struct {
	repo repository.FavoriteRepository
}

func NewService(repo repository.FavoriteRepository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Create(userID, favoriteType, bookID, seriesID string) (domain.Favorite, error) {
	item := domain.Favorite{
		ID:       idgen.NewFavorite(),
		UserID:   userID,
		Type:     favoriteType,
		BookID:   bookID,
		SeriesID: seriesID,
	}
	if err := s.repo.Create(item); err != nil {
		if errors.Is(err, repository.ErrFavoriteExists) {
			return domain.Favorite{}, ErrFavoriteExists
		}
		return domain.Favorite{}, err
	}
	return item, nil
}

func (s *Service) ListByUser(userID string) []domain.Favorite {
	return s.repo.ListByUser(userID)
}

func (s *Service) ListBySeriesID(seriesID string) []domain.Favorite {
	return s.repo.ListBySeriesID(seriesID)
}

func (s *Service) Delete(id string) bool {
	return s.repo.Delete(id)
}
