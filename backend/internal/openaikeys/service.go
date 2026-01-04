package openaikeys

import (
	"crypto/rand"
	"encoding/base64"
	"strings"
	"time"

	"book_manager/backend/internal/domain"
	"book_manager/backend/internal/repository"
)

type Service struct {
	repo repository.OpenAIKeyRepository
}

func NewService(repo repository.OpenAIKeyRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(name, apiKey string) domain.OpenAIKey {
	item := domain.OpenAIKey{
		ID:        newID(),
		Name:      strings.TrimSpace(name),
		APIKey:    strings.TrimSpace(apiKey),
		CreatedAt: time.Now(),
	}
	_ = s.repo.Create(item)
	return item
}

func (s *Service) List() []domain.OpenAIKey {
	return s.repo.List()
}

func (s *Service) Delete(id string) bool {
	return s.repo.Delete(id)
}

func (s *Service) First() (domain.OpenAIKey, bool) {
	items := s.repo.List()
	if len(items) == 0 {
		return domain.OpenAIKey{}, false
	}
	return items[0], true
}

func MaskKey(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	if len(value) <= 6 {
		return strings.Repeat("*", len(value))
	}
	head := value[:3]
	tail := value[len(value)-3:]
	return head + strings.Repeat("*", len(value)-6) + tail
}

func newID() string {
	seed := make([]byte, 16)
	if _, err := rand.Read(seed); err != nil {
		return "openai_key_demo"
	}
	return base64.RawURLEncoding.EncodeToString(seed)
}
