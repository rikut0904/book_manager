package repository

import "book_manager/backend/internal/domain"

type BookTagRepository interface {
	Create(item domain.BookTag) error
	Delete(item domain.BookTag) bool
	ListByUser(userID string) []domain.BookTag
	DeleteByTagID(tagID string)
}
