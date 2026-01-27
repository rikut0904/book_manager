package repository

import (
	"sync"

	"book_manager/backend/internal/domain"
)

type MemoryAdminInvitationRepository struct {
	mu       sync.RWMutex
	byID     map[string]domain.AdminInvitation
	byToken  map[string]string
	byUserID map[string]string
	ordered  []string
}

func NewMemoryAdminInvitationRepository() *MemoryAdminInvitationRepository {
	return &MemoryAdminInvitationRepository{
		byID:     make(map[string]domain.AdminInvitation),
		byToken:  make(map[string]string),
		byUserID: make(map[string]string),
	}
}

func (r *MemoryAdminInvitationRepository) Create(invitation domain.AdminInvitation) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.byID[invitation.ID] = invitation
	r.byToken[invitation.Token] = invitation.ID
	r.byUserID[invitation.UserID] = invitation.ID
	r.ordered = append(r.ordered, invitation.ID)
	return nil
}

func (r *MemoryAdminInvitationRepository) FindByID(id string) (domain.AdminInvitation, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	invitation, ok := r.byID[id]
	return invitation, ok
}

func (r *MemoryAdminInvitationRepository) FindByToken(token string) (domain.AdminInvitation, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	id, ok := r.byToken[token]
	if !ok {
		return domain.AdminInvitation{}, false
	}
	invitation, ok := r.byID[id]
	return invitation, ok
}

func (r *MemoryAdminInvitationRepository) FindByUserID(userID string) (domain.AdminInvitation, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	id, ok := r.byUserID[userID]
	if !ok {
		return domain.AdminInvitation{}, false
	}
	invitation, ok := r.byID[id]
	return invitation, ok
}

func (r *MemoryAdminInvitationRepository) List() []domain.AdminInvitation {
	r.mu.RLock()
	defer r.mu.RUnlock()

	items := make([]domain.AdminInvitation, 0, len(r.ordered))
	for _, id := range r.ordered {
		if invitation, ok := r.byID[id]; ok {
			items = append(items, invitation)
		}
	}
	return items
}

func (r *MemoryAdminInvitationRepository) Update(invitation domain.AdminInvitation) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.byID[invitation.ID]; !ok {
		return false
	}
	r.byID[invitation.ID] = invitation
	return true
}

func (r *MemoryAdminInvitationRepository) Delete(id string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	invitation, ok := r.byID[id]
	if !ok {
		return false
	}
	delete(r.byID, id)
	delete(r.byToken, invitation.Token)
	delete(r.byUserID, invitation.UserID)
	for i, storedID := range r.ordered {
		if storedID == id {
			r.ordered = append(r.ordered[:i], r.ordered[i+1:]...)
			break
		}
	}
	return true
}
