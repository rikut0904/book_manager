package gormrepo

import (
	"book_manager/backend/internal/domain"
	"book_manager/backend/internal/repository"
	"gorm.io/gorm"
)

type AdminInvitationRepository struct {
	db *gorm.DB
}

func NewAdminInvitationRepository(db *gorm.DB) *AdminInvitationRepository {
	return &AdminInvitationRepository{db: db}
}

func (r *AdminInvitationRepository) Create(invitation domain.AdminInvitation) error {
	model := AdminInvitation{
		ID:        invitation.ID,
		Token:     invitation.Token,
		Email:     invitation.Email,
		UserID:    invitation.UserID,
		CreatedBy: invitation.CreatedBy,
		ExpiresAt: invitation.ExpiresAt,
		UsedAt:    invitation.UsedAt,
		UsedBy:    invitation.UsedBy,
		CreatedAt: invitation.CreatedAt,
	}
	return r.db.Create(&model).Error
}

func (r *AdminInvitationRepository) FindByID(id string) (domain.AdminInvitation, bool) {
	var model AdminInvitation
	if err := r.db.Where("id = ?", id).First(&model).Error; err != nil {
		return domain.AdminInvitation{}, false
	}
	return toDomainInvitation(model), true
}

func (r *AdminInvitationRepository) FindByToken(token string) (domain.AdminInvitation, bool) {
	var model AdminInvitation
	if err := r.db.Where("token = ?", token).First(&model).Error; err != nil {
		return domain.AdminInvitation{}, false
	}
	return toDomainInvitation(model), true
}

func (r *AdminInvitationRepository) FindByUserID(userID string) (domain.AdminInvitation, bool) {
	var model AdminInvitation
	if err := r.db.Where("user_id = ?", userID).First(&model).Error; err != nil {
		return domain.AdminInvitation{}, false
	}
	return toDomainInvitation(model), true
}

func (r *AdminInvitationRepository) List() []domain.AdminInvitation {
	var models []AdminInvitation
	if err := r.db.Order("created_at desc").Find(&models).Error; err != nil {
		return nil
	}
	items := make([]domain.AdminInvitation, 0, len(models))
	for _, model := range models {
		items = append(items, toDomainInvitation(model))
	}
	return items
}

func (r *AdminInvitationRepository) Update(invitation domain.AdminInvitation) bool {
	result := r.db.Model(&AdminInvitation{}).Where("id = ?", invitation.ID).Updates(map[string]interface{}{
		"token":      invitation.Token,
		"email":      invitation.Email,
		"user_id":    invitation.UserID,
		"created_by": invitation.CreatedBy,
		"expires_at": invitation.ExpiresAt,
		"used_at":    invitation.UsedAt,
		"used_by":    invitation.UsedBy,
	})
	return result.Error == nil && result.RowsAffected > 0
}

func (r *AdminInvitationRepository) Delete(id string) bool {
	result := r.db.Delete(&AdminInvitation{}, "id = ?", id)
	return result.Error == nil && result.RowsAffected > 0
}

func toDomainInvitation(model AdminInvitation) domain.AdminInvitation {
	return domain.AdminInvitation{
		ID:        model.ID,
		Token:     model.Token,
		Email:     model.Email,
		UserID:    model.UserID,
		CreatedBy: model.CreatedBy,
		ExpiresAt: model.ExpiresAt,
		UsedAt:    model.UsedAt,
		UsedBy:    model.UsedBy,
		CreatedAt: model.CreatedAt,
	}
}

var _ repository.AdminInvitationRepository = (*AdminInvitationRepository)(nil)
