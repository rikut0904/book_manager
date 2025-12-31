package repository

import "book_manager/backend/internal/domain"

type TagRepository interface {
	Create(tag domain.Tag) error
	ListByUser(userID string) []domain.Tag
	FindByID(id string) (domain.Tag, bool)
	Delete(id string) bool
}
