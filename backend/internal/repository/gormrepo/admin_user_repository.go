package gormrepo

import (
	"book_manager/backend/internal/domain"
	"book_manager/backend/internal/repository"
	"gorm.io/gorm"
)

type AdminUserRepository struct {
	db *gorm.DB
}

func NewAdminUserRepository(db *gorm.DB) *AdminUserRepository {
	return &AdminUserRepository{db: db}
}

func (r *AdminUserRepository) Create(adminUser domain.AdminUser) error {
	model := AdminUser{
		ID:        adminUser.ID,
		UserID:    adminUser.UserID,
		CreatedBy: adminUser.CreatedBy,
		CreatedAt: adminUser.CreatedAt,
	}
	return r.db.Create(&model).Error
}

func (r *AdminUserRepository) FindByUserID(userID string) (domain.AdminUser, bool) {
	var model AdminUser
	if err := r.db.Where("user_id = ?", userID).First(&model).Error; err != nil {
		return domain.AdminUser{}, false
	}
	return toDomainAdminUser(model), true
}

func (r *AdminUserRepository) List() []domain.AdminUser {
	var models []AdminUser
	if err := r.db.Order("created_at asc").Find(&models).Error; err != nil {
		return nil
	}
	items := make([]domain.AdminUser, 0, len(models))
	for _, model := range models {
		items = append(items, toDomainAdminUser(model))
	}
	return items
}

func (r *AdminUserRepository) Delete(userID string) bool {
	result := r.db.Delete(&AdminUser{}, "user_id = ?", userID)
	return result.Error == nil && result.RowsAffected > 0
}

func toDomainAdminUser(model AdminUser) domain.AdminUser {
	return domain.AdminUser{
		ID:        model.ID,
		UserID:    model.UserID,
		CreatedBy: model.CreatedBy,
		CreatedAt: model.CreatedAt,
	}
}

var _ repository.AdminUserRepository = (*AdminUserRepository)(nil)
