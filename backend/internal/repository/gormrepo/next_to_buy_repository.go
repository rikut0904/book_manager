package gormrepo

import (
	"book_manager/backend/internal/domain"
	"book_manager/backend/internal/repository"
	"gorm.io/gorm"
)

type NextToBuyRepository struct {
	db *gorm.DB
}

func NewNextToBuyRepository(db *gorm.DB) *NextToBuyRepository {
	return &NextToBuyRepository{db: db}
}

func (r *NextToBuyRepository) Create(item domain.NextToBuyManual) error {
	model := NextToBuyManual{
		ID:           item.ID,
		UserID:       item.UserID,
		Title:        item.Title,
		SeriesName:   item.SeriesName,
		VolumeNumber: valueOrNilInt(item.VolumeNumber),
		Note:         item.Note,
	}
	if err := r.db.Create(&model).Error; err != nil {
		return err
	}
	return nil
}

func (r *NextToBuyRepository) ListByUser(userID string) []domain.NextToBuyManual {
	var models []NextToBuyManual
	if err := r.db.Where("user_id = ?", userID).Order("id asc").Find(&models).Error; err != nil {
		return nil
	}
	items := make([]domain.NextToBuyManual, 0, len(models))
	for _, model := range models {
		items = append(items, domain.NextToBuyManual{
			ID:           model.ID,
			UserID:       model.UserID,
			Title:        model.Title,
			SeriesName:   model.SeriesName,
			VolumeNumber: valueOrZeroInt(model.VolumeNumber),
			Note:         model.Note,
		})
	}
	return items
}

func (r *NextToBuyRepository) FindByID(id string) (domain.NextToBuyManual, bool) {
	var model NextToBuyManual
	if err := r.db.First(&model, "id = ?", id).Error; err != nil {
		return domain.NextToBuyManual{}, false
	}
	return domain.NextToBuyManual{
		ID:           model.ID,
		UserID:       model.UserID,
		Title:        model.Title,
		SeriesName:   model.SeriesName,
		VolumeNumber: valueOrZeroInt(model.VolumeNumber),
		Note:         model.Note,
	}, true
}

func (r *NextToBuyRepository) Update(item domain.NextToBuyManual) bool {
	model := NextToBuyManual{
		ID:           item.ID,
		UserID:       item.UserID,
		Title:        item.Title,
		SeriesName:   item.SeriesName,
		VolumeNumber: valueOrNilInt(item.VolumeNumber),
		Note:         item.Note,
	}
	if err := r.db.Save(&model).Error; err != nil {
		return false
	}
	return true
}

func (r *NextToBuyRepository) Delete(id string) bool {
	if err := r.db.Delete(&NextToBuyManual{}, "id = ?", id).Error; err != nil {
		return false
	}
	return true
}

var _ repository.NextToBuyRepository = (*NextToBuyRepository)(nil)
