package users

import (
	"errors"
	"strings"

	"book_manager/backend/internal/domain"
	"book_manager/backend/internal/repository"
)

var ErrInvalidVisibility = errors.New("invalid visibility")
var ErrUserNotFound = errors.New("user not found")
var ErrEmailExists = errors.New("email already exists")
var ErrUpdateFailed = errors.New("user update failed")

type Service struct {
	users    repository.UserRepository
	settings repository.ProfileSettingsRepository
}

func NewService(users repository.UserRepository, settings repository.ProfileSettingsRepository) *Service {
	return &Service{
		users:    users,
		settings: settings,
	}
}

func (s *Service) List(query string) []domain.User {
	items := s.users.List()
	if strings.TrimSpace(query) == "" {
		return items
	}
	normalized := strings.ToLower(strings.TrimSpace(query))
	result := make([]domain.User, 0, len(items))
	for _, user := range items {
		if strings.Contains(strings.ToLower(user.UserID), normalized) ||
			strings.Contains(strings.ToLower(user.Email), normalized) {
			result = append(result, user)
		}
	}
	return result
}

func (s *Service) Get(userID string) (domain.User, bool) {
	return s.users.FindByID(userID)
}

func (s *Service) Ensure(id, email, userID string) (domain.User, error) {
	if user, ok := s.users.FindByID(id); ok {
		return user, nil
	}
	user := domain.User{
		ID:          id,
		Email:       email,
		UserID:      userID,
		DisplayName: userID,
	}
	if err := s.users.Create(user); err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (s *Service) UpdateProfile(id string, displayName *string, email *string) (domain.User, error) {
	user, ok := s.users.FindByID(id)
	if !ok {
		return domain.User{}, ErrUserNotFound
	}
	if email != nil && *email != user.Email {
		if existing, ok := s.users.FindByEmail(*email); ok && existing.ID != id {
			return domain.User{}, ErrEmailExists
		}
		user.Email = *email
	}
	if displayName != nil {
		user.DisplayName = *displayName
	}
	if !s.users.Update(user) {
		return domain.User{}, ErrUpdateFailed
	}
	return user, nil
}

func (s *Service) Create(id, email, userID, displayName string) (domain.User, error) {
	user := domain.User{
		ID:          id,
		Email:       email,
		UserID:      userID,
		DisplayName: displayName,
	}
	if err := s.users.Create(user); err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (s *Service) IsUserIDTaken(userID string) bool {
	if strings.TrimSpace(userID) == "" {
		return false
	}
	_, ok := s.users.FindByUserID(userID)
	return ok
}

func (s *Service) Delete(userID string) bool {
	return s.users.Delete(userID)
}

func (s *Service) GetSettings(userID string) domain.ProfileSettings {
	if settings, ok := s.settings.Get(userID); ok {
		return settings
	}
	return domain.ProfileSettings{
		UserID:        userID,
		Visibility:    "public",
		OpenAIEnabled: false,
		OpenAIModel:   "",
		OpenAIAPIKey:  "",
	}
}

func (s *Service) UpdateSettings(userID, visibility string, openAIEnabled *bool, openAIModel string, openAIAPIKey string) (domain.ProfileSettings, error) {
	if visibility != "" && visibility != "public" && visibility != "followers" {
		return domain.ProfileSettings{}, ErrInvalidVisibility
	}
	current := s.GetSettings(userID)
	if visibility != "" {
		current.Visibility = visibility
	}
	if openAIEnabled != nil {
		current.OpenAIEnabled = *openAIEnabled
	}
	if openAIModel != "" {
		current.OpenAIModel = openAIModel
	}
	if openAIAPIKey != "" {
		current.OpenAIAPIKey = openAIAPIKey
	}
	s.settings.Upsert(current)
	return current, nil
}
