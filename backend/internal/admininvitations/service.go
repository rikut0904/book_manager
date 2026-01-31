package admininvitations

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"book_manager/backend/internal/domain"
	"book_manager/backend/internal/repository"

	"github.com/google/uuid"
)

var (
	ErrInvitationNotFound   = errors.New("invitation not found")
	ErrInvitationExpired    = errors.New("invitation expired")
	ErrInvitationUsed       = errors.New("invitation already used")
	ErrUserIDNotAdmin       = errors.New("user id is not an admin user id")
	ErrUserIDAlreadyTaken   = errors.New("user id already taken")
	ErrUserIDAlreadyInvited = errors.New("user id already has pending invitation")
	ErrEmailMismatch        = errors.New("email does not match invitation")
)

const (
	TokenBytes      = 32 // 256 bits
	ExpirationHours = 72
)

type Service struct {
	repo          repository.AdminInvitationRepository
	isAdminUserID func(userID string) bool
	isUserIDTaken func(userID string) bool
}

func NewService(
	repo repository.AdminInvitationRepository,
	isAdminUserID func(userID string) bool,
	isUserIDTaken func(userID string) bool,
) *Service {
	return &Service{
		repo:          repo,
		isAdminUserID: isAdminUserID,
		isUserIDTaken: isUserIDTaken,
	}
}

func (s *Service) Create(createdBy, userID, email string) (domain.AdminInvitation, error) {
	if s.isAdminUserID == nil || !s.isAdminUserID(userID) {
		return domain.AdminInvitation{}, ErrUserIDNotAdmin
	}
	if s.isUserIDTaken(userID) {
		return domain.AdminInvitation{}, ErrUserIDAlreadyTaken
	}
	if existing, ok := s.repo.FindByUserID(userID); ok && existing.UsedAt == nil {
		return domain.AdminInvitation{}, ErrUserIDAlreadyInvited
	}
	token, err := generateToken()
	if err != nil {
		return domain.AdminInvitation{}, err
	}
	now := time.Now()
	invitation := domain.AdminInvitation{
		ID:        uuid.New().String(),
		Token:     token,
		Email:     email,
		UserID:    userID,
		CreatedBy: createdBy,
		ExpiresAt: now.Add(ExpirationHours * time.Hour),
		CreatedAt: now,
	}
	if err := s.repo.Create(invitation); err != nil {
		return domain.AdminInvitation{}, err
	}
	return invitation, nil
}

func (s *Service) List() []domain.AdminInvitation {
	return s.repo.List()
}

func (s *Service) Get(id string) (domain.AdminInvitation, bool) {
	return s.repo.FindByID(id)
}

func (s *Service) Delete(id string) bool {
	return s.repo.Delete(id)
}

func (s *Service) ValidateToken(token string) (domain.AdminInvitation, error) {
	invitation, ok := s.repo.FindByToken(token)
	if !ok {
		return domain.AdminInvitation{}, ErrInvitationNotFound
	}
	if invitation.UsedAt != nil {
		return domain.AdminInvitation{}, ErrInvitationUsed
	}
	if time.Now().After(invitation.ExpiresAt) {
		return domain.AdminInvitation{}, ErrInvitationExpired
	}
	return invitation, nil
}

func (s *Service) ValidateAndCheckEmail(token, email string) (domain.AdminInvitation, error) {
	invitation, err := s.ValidateToken(token)
	if err != nil {
		return domain.AdminInvitation{}, err
	}
	if invitation.Email != "" && invitation.Email != email {
		return domain.AdminInvitation{}, ErrEmailMismatch
	}
	return invitation, nil
}

func (s *Service) MarkUsed(id, usedBy string) error {
	invitation, ok := s.repo.FindByID(id)
	if !ok {
		return ErrInvitationNotFound
	}
	now := time.Now()
	invitation.UsedAt = &now
	invitation.UsedBy = usedBy
	if !s.repo.Update(invitation) {
		return ErrInvitationNotFound
	}
	return nil
}

// UnmarkUsed はロールバック用に招待を未使用状態に戻します
func (s *Service) UnmarkUsed(id string) error {
	invitation, ok := s.repo.FindByID(id)
	if !ok {
		return ErrInvitationNotFound
	}
	invitation.UsedAt = nil
	invitation.UsedBy = ""
	if !s.repo.Update(invitation) {
		return ErrInvitationNotFound
	}
	return nil
}

func (s *Service) IsAdminUserID(userID string) bool {
	if s.isAdminUserID == nil {
		return false
	}
	return s.isAdminUserID(userID)
}

func generateToken() (string, error) {
	bytes := make([]byte, TokenBytes)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
