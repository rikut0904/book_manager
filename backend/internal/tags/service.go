package tags

import (
	"crypto/rand"
	"encoding/base64"
	"errors"

	"book_manager/backend/internal/domain"
	"book_manager/backend/internal/repository"
)

var (
	ErrTagExists     = errors.New("tag already exists")
	ErrTagNotFound   = errors.New("tag not found")
	ErrBookTagExists = errors.New("book tag already exists")
)

type Service struct {
	tags     repository.TagRepository
	bookTags repository.BookTagRepository
}

func NewService(tags repository.TagRepository, bookTags repository.BookTagRepository) *Service {
	return &Service{
		tags:     tags,
		bookTags: bookTags,
	}
}

func (s *Service) Create(userID, name string) (domain.Tag, error) {
	tag := domain.Tag{
		ID:          newID(),
		OwnerUserID: userID,
		Name:        name,
	}
	if err := s.tags.Create(tag); err != nil {
		if errors.Is(err, repository.ErrTagExists) {
			return domain.Tag{}, ErrTagExists
		}
		return domain.Tag{}, err
	}
	return tag, nil
}

func (s *Service) ListByUser(userID string) []domain.Tag {
	return s.tags.ListByUser(userID)
}

func (s *Service) Delete(tagID string) bool {
	ok := s.tags.Delete(tagID)
	if ok {
		s.bookTags.DeleteByTagID(tagID)
	}
	return ok
}

func (s *Service) AddBookTag(userID, bookID, tagID string) (domain.BookTag, error) {
	if _, ok := s.tags.FindByID(tagID); !ok {
		return domain.BookTag{}, ErrTagNotFound
	}
	item := domain.BookTag{
		UserID: userID,
		BookID: bookID,
		TagID:  tagID,
	}
	if err := s.bookTags.Create(item); err != nil {
		if errors.Is(err, repository.ErrBookTagExists) {
			return domain.BookTag{}, ErrBookTagExists
		}
		return domain.BookTag{}, err
	}
	return item, nil
}

func (s *Service) RemoveBookTag(userID, bookID, tagID string) bool {
	return s.bookTags.Delete(domain.BookTag{
		UserID: userID,
		BookID: bookID,
		TagID:  tagID,
	})
}

func newID() string {
	seed := make([]byte, 16)
	if _, err := rand.Read(seed); err != nil {
		return "tag_demo"
	}
	return base64.RawURLEncoding.EncodeToString(seed)
}
