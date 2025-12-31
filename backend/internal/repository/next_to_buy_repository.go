package repository

import "book_manager/backend/internal/domain"

type NextToBuyRepository interface {
	Create(item domain.NextToBuyManual) error
	ListByUser(userID string) []domain.NextToBuyManual
	FindByID(id string) (domain.NextToBuyManual, bool)
	Update(item domain.NextToBuyManual) bool
	Delete(id string) bool
}
