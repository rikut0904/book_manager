package reports

import (
	"fmt"
	"log"
	"net/smtp"
	"strings"

	"book_manager/backend/internal/domain"
)

type Service struct {
	to   string
	smtp SMTPConfig
}

type SMTPConfig struct {
	Host string
	Port string
	User string
	Pass string
	From string
}

func NewService(to string, smtpConfig SMTPConfig) *Service {
	return &Service{
		to:   to,
		smtp: smtpConfig,
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
		"ISBN: %s\nTitle: %s\nAuthors: %v\nPublisher: %s\nPublishedDate: %s\nThumbnail: %s\n\nSuggestion:\n%s\n\nNote:\n%s\n",
		report.Book.ISBN13,
		report.Book.Title,
		report.Book.Authors,
		report.Book.Publisher,
		report.Book.PublishedDate,
		report.Book.ThumbnailURL,
		report.Suggestion,
		report.Note,
	)

	if s.smtp.Host == "" {
		log.Printf("book report email (dry-run):\nTo: %s\nSubject: %s\n\n%s", s.to, subject, body)
		return
	}

	from := s.smtp.From
	if strings.TrimSpace(from) == "" {
		from = s.smtp.User
	}
	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s", from, s.to, subject, body)
	addr := fmt.Sprintf("%s:%s", s.smtp.Host, s.smtp.Port)
	var auth smtp.Auth
	if s.smtp.User != "" && s.smtp.Pass != "" {
		auth = smtp.PlainAuth("", s.smtp.User, s.smtp.Pass, s.smtp.Host)
	}
	if err := smtp.SendMail(addr, auth, from, []string{s.to}, []byte(msg)); err != nil {
		log.Printf("book report email send error: %v", err)
	}
}
