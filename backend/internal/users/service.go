package users

import (
	"errors"
	"strings"

	"book_manager/backend/internal/domain"
	"book_manager/backend/internal/repository"
)

var ErrInvalidVisibility = errors.New("invalid visibility")

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
		if strings.Contains(strings.ToLower(user.Username), normalized) ||
			strings.Contains(strings.ToLower(user.Email), normalized) {
			result = append(result, user)
		}
	}
	return result
}

func (s *Service) Get(userID string) (domain.User, bool) {
	return s.users.FindByID(userID)
}

func (s *Service) UpdateUsername(userID, username string) (domain.User, bool) {
	user, ok := s.users.FindByID(userID)
	if !ok {
		return domain.User{}, false
	}
	user.Username = username
	if !s.users.Update(user) {
		return domain.User{}, false
	}
	return user, true
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
		GeminiEnabled: false,
		GeminiModel:   "",
		GeminiAPIKey:  "",
	}
}

func (s *Service) UpdateSettings(userID, visibility string, geminiEnabled *bool, geminiModel string, geminiAPIKey string) (domain.ProfileSettings, error) {
	if visibility != "" && visibility != "public" && visibility != "followers" {
		return domain.ProfileSettings{}, ErrInvalidVisibility
	}
	current := s.GetSettings(userID)
	if visibility != "" {
		current.Visibility = visibility
	}
	if geminiEnabled != nil {
		current.GeminiEnabled = *geminiEnabled
	}
	if geminiModel != "" {
		current.GeminiModel = geminiModel
	}
	if geminiAPIKey != "" {
		current.GeminiAPIKey = geminiAPIKey
	}
	s.settings.Upsert(current)
	return current, nil
}
