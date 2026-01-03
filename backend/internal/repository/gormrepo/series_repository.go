package gormrepo

import (
	"book_manager/backend/internal/domain"
	"book_manager/backend/internal/repository"
	"gorm.io/gorm"
)

type SeriesRepository struct {
	db *gorm.DB
}

func NewSeriesRepository(db *gorm.DB) *SeriesRepository {
	return &SeriesRepository{db: db}
}

func (r *SeriesRepository) Create(series domain.Series) error {
	model := Series{
		ID:             series.ID,
		Name:           series.Name,
		NormalizedName: series.NormalizedName,
	}
	if err := r.db.Create(&model).Error; err != nil {
		if isUniqueViolation(err) {
			return repository.ErrSeriesExists
		}
		return err
	}
	return nil
}

func (r *SeriesRepository) List() []domain.Series {
	var models []Series
	if err := r.db.Order("name asc").Find(&models).Error; err != nil {
		return nil
	}
	items := make([]domain.Series, 0, len(models))
	for _, model := range models {
		items = append(items, domain.Series{
			ID:             model.ID,
			Name:           model.Name,
			NormalizedName: model.NormalizedName,
		})
	}
	return items
}

func (r *SeriesRepository) FindByNormalizedName(name string) (domain.Series, bool) {
	var model Series
	if err := r.db.First(&model, "normalized_name = ?", name).Error; err != nil {
		return domain.Series{}, false
	}
	return domain.Series{
		ID:             model.ID,
		Name:           model.Name,
		NormalizedName: model.NormalizedName,
	}, true
}

func (r *SeriesRepository) Delete(id string) bool {
	result := r.db.Delete(&Series{}, "id = ?", id)
	if result.Error != nil {
		return false
	}
	return result.RowsAffected > 0
}

func (r *SeriesRepository) Update(series domain.Series) bool {
	model := Series{
		ID:             series.ID,
		Name:           series.Name,
		NormalizedName: series.NormalizedName,
	}
	if err := r.db.Model(&Series{}).Where("id = ?", series.ID).Updates(&model).Error; err != nil {
		return false
	}
	return true
}

var _ repository.SeriesRepository = (*SeriesRepository)(nil)
