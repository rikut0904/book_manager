package repository

import "book_manager/backend/internal/domain"

type OpenAIKeyRepository interface {
	Create(key domain.OpenAIKey) error
	List() []domain.OpenAIKey
	Delete(id string) bool
}
