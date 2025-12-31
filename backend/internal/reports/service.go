package reports

import (
	"fmt"
	"log"

	"book_manager/backend/internal/domain"
)

type Service struct {
	to string
}

func NewService(to string) *Service {
	return &Service{
		to: to,
	}
}

type BookReport struct {
	BookID     string
	Suggestion string
	Note       string
	Book       domain.Book
}

func (s *Service) SendBookReport(report BookReport) {
	subject := fmt.Sprintf("[BookManager] 書誌情報の修正提案 (%s)", report.BookID)
	body := fmt.Sprintf(
		"To: %s\nSubject: %s\n\nISBN: %s\nTitle: %s\nAuthors: %v\nPublisher: %s\nPublishedDate: %s\nThumbnail: %s\n\nSuggestion:\n%s\n\nNote:\n%s\n",
		s.to,
		subject,
		report.Book.ISBN13,
		report.Book.Title,
		report.Book.Authors,
		report.Book.Publisher,
		report.Book.PublishedDate,
		report.Book.ThumbnailURL,
		report.Suggestion,
		report.Note,
	)
	log.Printf("book report email:\n%s", body)
}
