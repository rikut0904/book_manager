package gormrepo

import (
	"book_manager/backend/internal/domain"
	"book_manager/backend/internal/repository"
	"gorm.io/gorm"
)

type TagRepository struct {
	db *gorm.DB
}

func NewTagRepository(db *gorm.DB) *TagRepository {
	return &TagRepository{db: db}
}

func (r *TagRepository) Create(tag domain.Tag) error {
	model := Tag{
		ID:          tag.ID,
		OwnerUserID: tag.OwnerUserID,
		Name:        tag.Name,
	}
	if err := r.db.Create(&model).Error; err != nil {
		if isUniqueViolation(err) {
			return repository.ErrTagExists
		}
		return err
	}
	return nil
}

func (r *TagRepository) ListByUser(userID string) []domain.Tag {
	var models []Tag
	if err := r.db.Where("owner_user_id = ?", userID).Order("id asc").Find(&models).Error; err != nil {
		return nil
	}
	items := make([]domain.Tag, 0, len(models))
	for _, model := range models {
		items = append(items, domain.Tag{
			ID:          model.ID,
			OwnerUserID: model.OwnerUserID,
			Name:        model.Name,
		})
	}
	return items
}

func (r *TagRepository) FindByID(id string) (domain.Tag, bool) {
	var model Tag
	if err := r.db.First(&model, "id = ?", id).Error; err != nil {
		return domain.Tag{}, false
	}
	return domain.Tag{
		ID:          model.ID,
		OwnerUserID: model.OwnerUserID,
		Name:        model.Name,
	}, true
}

func (r *TagRepository) Delete(id string) bool {
	if err := r.db.Delete(&Tag{}, "id = ?", id).Error; err != nil {
		return false
	}
	return true
}

var _ repository.TagRepository = (*TagRepository)(nil)
