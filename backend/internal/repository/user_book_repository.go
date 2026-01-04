package repository

import "book_manager/backend/internal/domain"

type UserBookRepository interface {
	Create(userBook domain.UserBook) error
	ListByUser(userID string) []domain.UserBook
	ListAll() []domain.UserBook
	ListBySeriesID(seriesID string) []domain.UserBook
	FindByID(id string) (domain.UserBook, bool)
	Update(userBook domain.UserBook) bool
	Delete(id string) bool
}
