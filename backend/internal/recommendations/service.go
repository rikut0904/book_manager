package recommendations

import (
	"crypto/rand"
	"encoding/base64"
	"time"

	"book_manager/backend/internal/domain"
	"book_manager/backend/internal/repository"
)

type Service struct {
	repo repository.RecommendationRepository
}

func NewService(repo repository.RecommendationRepository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Create(userID, bookID, comment string) (domain.Recommendation, error) {
	item := domain.Recommendation{
		ID:        newID(),
		UserID:    userID,
		BookID:    bookID,
		Comment:   comment,
		CreatedAt: time.Now(),
	}
	if err := s.repo.Create(item); err != nil {
		return domain.Recommendation{}, err
	}
	return item, nil
}

func (s *Service) List() []domain.Recommendation {
	return s.repo.List()
}

func (s *Service) Delete(id string) bool {
	return s.repo.Delete(id)
}

func newID() string {
	seed := make([]byte, 16)
	if _, err := rand.Read(seed); err != nil {
		return "recommendation_demo"
	}
	return base64.RawURLEncoding.EncodeToString(seed)
}
