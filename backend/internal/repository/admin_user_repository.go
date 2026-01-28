package repository

import "book_manager/backend/internal/domain"

type AdminUserRepository interface {
	Create(adminUser domain.AdminUser) error
	FindByUserID(userID string) (domain.AdminUser, bool)
	List() []domain.AdminUser
	Delete(userID string) bool
}
