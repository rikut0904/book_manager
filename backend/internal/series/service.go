package series

import (
	"errors"
	"strings"

	"book_manager/backend/internal/domain"
	"book_manager/backend/internal/idgen"
	"book_manager/backend/internal/isbn"
	"book_manager/backend/internal/repository"
)

var ErrSeriesExists = errors.New("series already exists")

type Service struct {
	repo repository.SeriesRepository
}

func NewService(repo repository.SeriesRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(name string) (domain.Series, error) {
	cleaned := sanitizeName(name)
	normalized := normalizeName(cleaned)
	series := domain.Series{
		ID:             idgen.NewSeries(),
		Name:           cleaned,
		NormalizedName: normalized,
	}
	if err := s.repo.Create(series); err != nil {
		if errors.Is(err, repository.ErrSeriesExists) {
			return domain.Series{}, ErrSeriesExists
		}
		return domain.Series{}, err
	}
	return series, nil
}

func (s *Service) Ensure(name string) (domain.Series, error) {
	cleaned := sanitizeName(name)
	normalized := normalizeName(cleaned)
	if item, ok := s.repo.FindByNormalizedName(normalized); ok {
		return item, nil
	}
	return s.Create(cleaned)
}

func (s *Service) List() []domain.Series {
	return s.repo.List()
}

func (s *Service) Delete(id string) bool {
	return s.repo.Delete(id)
}

func (s *Service) NormalizeAll() int {
	items := s.repo.List()
	updated := 0
	for _, item := range items {
		cleaned := sanitizeName(item.Name)
		if cleaned == "" {
			continue
		}
		normalized := normalizeName(cleaned)
		if item.Name == cleaned && item.NormalizedName == normalized {
			continue
		}
		if existing, ok := s.repo.FindByNormalizedName(normalized); ok && existing.ID != item.ID {
			continue
		}
		item.Name = cleaned
		item.NormalizedName = normalized
		if s.repo.Update(item) {
			updated++
		}
	}
	return updated
}

func normalizeName(name string) string {
	return strings.ToLower(strings.TrimSpace(name))
}

func sanitizeName(name string) string {
	cleaned := isbn.NormalizeSeriesName(name)
	if cleaned == "" {
		return strings.TrimSpace(name)
	}
	return cleaned
}
