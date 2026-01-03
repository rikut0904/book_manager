package repository

import (
	"errors"

	"book_manager/backend/internal/domain"
)

var ErrSeriesExists = errors.New("series already exists")

type SeriesRepository interface {
	Create(series domain.Series) error
	List() []domain.Series
	FindByNormalizedName(name string) (domain.Series, bool)
	Delete(id string) bool
	Update(series domain.Series) bool
}
