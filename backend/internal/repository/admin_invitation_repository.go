package repository

import "book_manager/backend/internal/domain"

type AdminInvitationRepository interface {
	Create(invitation domain.AdminInvitation) error
	FindByID(id string) (domain.AdminInvitation, bool)
	FindByToken(token string) (domain.AdminInvitation, bool)
	FindByUserID(userID string) (domain.AdminInvitation, bool)
	List() []domain.AdminInvitation
	Update(invitation domain.AdminInvitation) bool
	Delete(id string) bool
}
