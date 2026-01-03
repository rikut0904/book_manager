package gormrepo

import (
	"book_manager/backend/internal/domain"
	"book_manager/backend/internal/repository"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type BookRepository struct {
	db *gorm.DB
}

func NewBookRepository(db *gorm.DB) *BookRepository {
	return &BookRepository{db: db}
}

func (r *BookRepository) Create(book domain.Book) error {
	var isbnPtr *string
	if book.ISBN13 != "" {
		isbnPtr = &book.ISBN13
	}
	authors, err := marshalAuthors(book.Authors)
	if err != nil {
		return err
	}
	model := Book{
		ID:            book.ID,
		ISBN13:        isbnPtr,
		Title:         book.Title,
		Authors:       datatypes.JSON(authors),
		Publisher:     book.Publisher,
		PublishedDate: book.PublishedDate,
		ThumbnailURL:  book.ThumbnailURL,
		Source:        book.Source,
		SeriesName:    book.SeriesName,
	}
	if err := r.db.Create(&model).Error; err != nil {
		if isUniqueViolation(err) {
			return repository.ErrBookExists
		}
		return err
	}
	return nil
}

func (r *BookRepository) FindByID(id string) (domain.Book, bool) {
	var model Book
	if err := r.db.First(&model, "id = ?", id).Error; err != nil {
		return domain.Book{}, false
	}
	return modelToDomainBook(model), true
}

func (r *BookRepository) FindByISBN(isbn string) (domain.Book, bool) {
	var model Book
	if err := r.db.Where("isbn13 = ?", isbn).First(&model).Error; err != nil {
		return domain.Book{}, false
	}
	return modelToDomainBook(model), true
}

func (r *BookRepository) List() []domain.Book {
	var models []Book
	if err := r.db.Order("id asc").Find(&models).Error; err != nil {
		return nil
	}
	items := make([]domain.Book, 0, len(models))
	for _, model := range models {
		items = append(items, modelToDomainBook(model))
	}
	return items
}

func (r *BookRepository) Delete(id string) bool {
	result := r.db.Delete(&Book{}, "id = ?", id)
	if result.Error != nil {
		return false
	}
	return result.RowsAffected > 0
}

func modelToDomainBook(model Book) domain.Book {
	isbn := ""
	if model.ISBN13 != nil {
		isbn = *model.ISBN13
	}
	return domain.Book{
		ID:            model.ID,
		ISBN13:        isbn,
		Title:         model.Title,
		Authors:       unmarshalAuthors(model.Authors),
		Publisher:     model.Publisher,
		PublishedDate: model.PublishedDate,
		ThumbnailURL:  model.ThumbnailURL,
		Source:        model.Source,
		SeriesName:    model.SeriesName,
	}
}

var _ repository.BookRepository = (*BookRepository)(nil)
