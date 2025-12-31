package favorites

import (
	"crypto/rand"
	"encoding/base64"
	"errors"

	"book_manager/backend/internal/domain"
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
		ID:       newID(),
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

func (s *Service) Delete(id string) bool {
	return s.repo.Delete(id)
}

func newID() string {
	seed := make([]byte, 16)
	if _, err := rand.Read(seed); err != nil {
		return "favorite_demo"
	}
	return base64.RawURLEncoding.EncodeToString(seed)
}
