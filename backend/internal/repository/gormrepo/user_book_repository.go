package gormrepo

import (
	"book_manager/backend/internal/domain"
	"book_manager/backend/internal/repository"
	"gorm.io/gorm"
)

type UserBookRepository struct {
	db *gorm.DB
}

func NewUserBookRepository(db *gorm.DB) *UserBookRepository {
	return &UserBookRepository{db: db}
}

func (r *UserBookRepository) Create(userBook domain.UserBook) error {
	model := UserBook{
		ID:           userBook.ID,
		UserID:       userBook.UserID,
		BookID:       userBook.BookID,
		Note:         userBook.Note,
		AcquiredAt:   userBook.AcquiredAt,
		SeriesID:     userBook.SeriesID,
		VolumeNumber: valueOrNilInt(userBook.VolumeNumber),
		SeriesSource: userBook.SeriesSource,
	}
	if err := r.db.Create(&model).Error; err != nil {
		if isUniqueViolation(err) {
			return repository.ErrUserBookExists
		}
		return err
	}
	return nil
}

func (r *UserBookRepository) ListByUser(userID string) []domain.UserBook {
	var models []UserBook
	if err := r.db.Where("user_id = ?", userID).Order("id asc").Find(&models).Error; err != nil {
		return nil
	}
	items := make([]domain.UserBook, 0, len(models))
	for _, model := range models {
		items = append(items, modelToDomainUserBook(model))
	}
	return items
}

func (r *UserBookRepository) ListAll() []domain.UserBook {
	var models []UserBook
	if err := r.db.Order("id asc").Find(&models).Error; err != nil {
		return nil
	}
	items := make([]domain.UserBook, 0, len(models))
	for _, model := range models {
		items = append(items, modelToDomainUserBook(model))
	}
	return items
}

func (r *UserBookRepository) ListBySeriesID(seriesID string) []domain.UserBook {
	var models []UserBook
	if err := r.db.Where("series_id = ?", seriesID).Order("id asc").Find(&models).Error; err != nil {
		return nil
	}
	items := make([]domain.UserBook, 0, len(models))
	for _, model := range models {
		items = append(items, modelToDomainUserBook(model))
	}
	return items
}

func (r *UserBookRepository) FindByID(id string) (domain.UserBook, bool) {
	var model UserBook
	if err := r.db.First(&model, "id = ?", id).Error; err != nil {
		return domain.UserBook{}, false
	}
	return modelToDomainUserBook(model), true
}

func (r *UserBookRepository) Update(userBook domain.UserBook) bool {
	model := UserBook{
		ID:           userBook.ID,
		UserID:       userBook.UserID,
		BookID:       userBook.BookID,
		Note:         userBook.Note,
		AcquiredAt:   userBook.AcquiredAt,
		SeriesID:     userBook.SeriesID,
		VolumeNumber: valueOrNilInt(userBook.VolumeNumber),
		SeriesSource: userBook.SeriesSource,
	}
	if err := r.db.Save(&model).Error; err != nil {
		return false
	}
	return true
}

func (r *UserBookRepository) Delete(id string) bool {
	if err := r.db.Delete(&UserBook{}, "id = ?", id).Error; err != nil {
		return false
	}
	return true
}

func modelToDomainUserBook(model UserBook) domain.UserBook {
	return domain.UserBook{
		ID:           model.ID,
		UserID:       model.UserID,
		BookID:       model.BookID,
		Note:         model.Note,
		AcquiredAt:   model.AcquiredAt,
		SeriesID:     model.SeriesID,
		VolumeNumber: valueOrZeroInt(model.VolumeNumber),
		SeriesSource: model.SeriesSource,
	}
}

func valueOrNilInt(value int) *int {
	if value == 0 {
		return nil
	}
	return &value
}

func valueOrZeroInt(value *int) int {
	if value == nil {
		return 0
	}
	return *value
}

var _ repository.UserBookRepository = (*UserBookRepository)(nil)
