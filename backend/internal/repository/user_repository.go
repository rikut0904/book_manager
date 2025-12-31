package repository

import "book_manager/backend/internal/domain"

type UserRepository interface {
	Create(user domain.User) error
	FindByEmail(email string) (domain.User, bool)
	FindByID(id string) (domain.User, bool)
}
