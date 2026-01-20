package gormrepo

import (
	"book_manager/backend/internal/domain"
	"book_manager/backend/internal/repository"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user domain.User) error {
	model := User{
		ID:           user.ID,
		Email:        user.Email,
		UserID:     user.UserID,
		DisplayName:  user.DisplayName,
		PasswordHash: user.PasswordHash,
	}
	if _, ok := r.FindByID(user.ID); ok {
		return repository.ErrUserExists
	}
	if err := r.db.Create(&model).Error; err != nil {
		if isUniqueViolation(err) {
			return repository.ErrUserExists
		}
		return err
	}
	return nil
}

func (r *UserRepository) FindByEmail(email string) (domain.User, bool) {
	var model User
	if err := r.db.Where("email = ?", email).First(&model).Error; err != nil {
		return domain.User{}, false
	}
	return domain.User{
		ID:           model.ID,
		Email:        model.Email,
		UserID:     model.UserID,
		DisplayName:  model.DisplayName,
		PasswordHash: model.PasswordHash,
	}, true
}

func (r *UserRepository) FindByUserID(userID string) (domain.User, bool) {
	var model User
	if err := r.db.Where("user_id = ?", userID).First(&model).Error; err != nil {
		return domain.User{}, false
	}
	return domain.User{
		ID:           model.ID,
		Email:        model.Email,
		UserID:     model.UserID,
		DisplayName:  model.DisplayName,
		PasswordHash: model.PasswordHash,
	}, true
}

func (r *UserRepository) FindByID(id string) (domain.User, bool) {
	var model User
	if err := r.db.First(&model, "id = ?", id).Error; err != nil {
		return domain.User{}, false
	}
	return domain.User{
		ID:           model.ID,
		Email:        model.Email,
		UserID:     model.UserID,
		DisplayName:  model.DisplayName,
		PasswordHash: model.PasswordHash,
	}, true
}

func (r *UserRepository) List() []domain.User {
	var models []User
	if err := r.db.Order("id asc").Find(&models).Error; err != nil {
		return nil
	}
	items := make([]domain.User, 0, len(models))
	for _, model := range models {
		items = append(items, domain.User{
			ID:           model.ID,
			Email:        model.Email,
			UserID:     model.UserID,
			DisplayName:  model.DisplayName,
			PasswordHash: model.PasswordHash,
		})
	}
	return items
}

func (r *UserRepository) Update(user domain.User) bool {
	model := User{
		ID:           user.ID,
		Email:        user.Email,
		UserID:     user.UserID,
		DisplayName:  user.DisplayName,
		PasswordHash: user.PasswordHash,
	}
	if err := r.db.Save(&model).Error; err != nil {
		return false
	}
	return true
}

func (r *UserRepository) Delete(id string) bool {
	if err := r.db.Delete(&User{}, "id = ?", id).Error; err != nil {
		return false
	}
	return true
}

var _ repository.UserRepository = (*UserRepository)(nil)
