package repository

import "book_manager/backend/internal/domain"

type FavoriteRepository interface {
	Create(favorite domain.Favorite) error
	ListByUser(userID string) []domain.Favorite
	ListBySeriesID(seriesID string) []domain.Favorite
	Delete(id string) bool
}
