package repository

import "book_manager/backend/internal/domain"

type RecommendationRepository interface {
	Create(item domain.Recommendation) error
	List() []domain.Recommendation
	ListByUser(userID string) []domain.Recommendation
	Delete(id string) bool
}
