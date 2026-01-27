package adminusers

import (
	"errors"
	"sort"
	"time"

	"book_manager/backend/internal/domain"
	"book_manager/backend/internal/repository"

	"github.com/google/uuid"
)

var (
	ErrAlreadyAdmin = errors.New("user is already an admin")
	ErrNotDBAdmin   = errors.New("user is not a DB admin")
	ErrEnvAdmin     = errors.New("cannot remove env admin from DB")
)

type Service struct {
	repo        repository.AdminUserRepository
	envAdminIDs map[string]struct{}
}

func NewService(repo repository.AdminUserRepository, envAdminIDs []string) *Service {
	envMap := make(map[string]struct{}, len(envAdminIDs))
	for _, userID := range envAdminIDs {
		if userID != "" {
			envMap[userID] = struct{}{}
		}
	}
	return &Service{
		repo:        repo,
		envAdminIDs: envMap,
	}
}

func (s *Service) IsAdmin(userID string) bool {
	if userID == "" {
		return false
	}
	if _, ok := s.envAdminIDs[userID]; ok {
		return true
	}
	_, ok := s.repo.FindByUserID(userID)
	return ok
}

func (s *Service) List() []AdminUserInfo {
	dbAdmins := s.repo.List()
	result := make([]AdminUserInfo, 0, len(s.envAdminIDs)+len(dbAdmins))

	dbUserIDs := make(map[string]struct{}, len(dbAdmins))
	for _, admin := range dbAdmins {
		dbUserIDs[admin.UserID] = struct{}{}
	}

	for userID := range s.envAdminIDs {
		_, inDB := dbUserIDs[userID]
		result = append(result, AdminUserInfo{
			UserID:    userID,
			Source:    "env",
			CreatedAt: time.Time{},
			InDB:      inDB,
		})
	}

	for _, admin := range dbAdmins {
		if _, isEnv := s.envAdminIDs[admin.UserID]; isEnv {
			continue
		}
		result = append(result, AdminUserInfo{
			UserID:    admin.UserID,
			Source:    "db",
			CreatedBy: admin.CreatedBy,
			CreatedAt: admin.CreatedAt,
			InDB:      true,
		})
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].UserID < result[j].UserID
	})
	return result
}

func (s *Service) ListAll() []string {
	dbAdmins := s.repo.List()
	seen := make(map[string]struct{}, len(s.envAdminIDs)+len(dbAdmins))
	result := make([]string, 0, len(s.envAdminIDs)+len(dbAdmins))

	for userID := range s.envAdminIDs {
		if _, ok := seen[userID]; !ok {
			seen[userID] = struct{}{}
			result = append(result, userID)
		}
	}
	for _, admin := range dbAdmins {
		if _, ok := seen[admin.UserID]; !ok {
			seen[admin.UserID] = struct{}{}
			result = append(result, admin.UserID)
		}
	}
	return result
}

func (s *Service) Add(userID, createdBy string) error {
	if _, ok := s.repo.FindByUserID(userID); ok {
		return ErrAlreadyAdmin
	}

	adminUser := domain.AdminUser{
		ID:        uuid.New().String(),
		UserID:    userID,
		CreatedBy: createdBy,
		CreatedAt: time.Now(),
	}
	return s.repo.Create(adminUser)
}

func (s *Service) Remove(userID string) error {
	if _, ok := s.envAdminIDs[userID]; ok {
		if _, inDB := s.repo.FindByUserID(userID); !inDB {
			return ErrEnvAdmin
		}
	}

	if !s.repo.Delete(userID) {
		return ErrNotDBAdmin
	}
	return nil
}

func (s *Service) IsEnvAdmin(userID string) bool {
	_, ok := s.envAdminIDs[userID]
	return ok
}

type AdminUserInfo struct {
	UserID    string    `json:"userId"`
	Source    string    `json:"source"`
	CreatedBy string    `json:"createdBy,omitempty"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
	InDB      bool      `json:"inDb"`
}
