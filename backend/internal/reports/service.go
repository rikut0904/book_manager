package reports

import (
	"bytes"
	"fmt"
	"log"
	"net/smtp"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"book_manager/backend/internal/domain"
)

type Service struct {
	to           string
	smtp         SMTPConfig
	templatesDir string
	frontendURL  string
}

type SMTPConfig struct {
	Host string
	Port string
	User string
	Pass string
	From string
}

func NewService(to string, smtpConfig SMTPConfig, templatesDir string, frontendURL string) *Service {
	return &Service{
		to:           to,
		smtp:         smtpConfig,
		templatesDir: templatesDir,
		frontendURL:  frontendURL,
	}
}

func (s *Service) loadTemplate(name string, data interface{}) (string, error) {
	path := filepath.Join(s.templatesDir, "email", name)
	content, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to read template %s: %w", name, err)
	}

	tmpl, err := template.New(name).Parse(string(content))
	if err != nil {
		return "", fmt.Errorf("failed to parse template %s: %w", name, err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template %s: %w", name, err)
	}

	return strings.TrimSpace(buf.String()), nil
}

type BookReport struct {
	BookID     string
	Suggestion string
	Note       string
	Book       domain.Book
}

func (s *Service) SendBookReport(report BookReport) {
	data := map[string]interface{}{
		"BookID":        report.BookID,
		"ISBN13":        report.Book.ISBN13,
		"Title":         report.Book.Title,
		"Authors":       strings.Join(report.Book.Authors, ", "),
		"Publisher":     report.Book.Publisher,
		"PublishedDate": report.Book.PublishedDate,
		"ThumbnailURL":  report.Book.ThumbnailURL,
		"Suggestion":    report.Suggestion,
		"Note":          report.Note,
	}

	subject, err := s.loadTemplate("book_report_subject.txt", data)
	if err != nil {
		log.Printf("failed to load book_report_subject template: %v", err)
		subject = fmt.Sprintf("[BookManager] 書誌情報の修正提案 (%s)", report.BookID)
	}

	body, err := s.loadTemplate("book_report_body.txt", data)
	if err != nil {
		log.Printf("failed to load book_report_body template: %v", err)
		body = fmt.Sprintf(
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
	}

	s.sendMail(s.to, subject, body)
}

type AdminInvitationEmail struct {
	To        string
	UserID    string
	Token     string
	ExpiresAt time.Time
}

func (s *Service) SendAdminInvitation(invitation AdminInvitationEmail) {
	if strings.TrimSpace(invitation.To) == "" {
		return
	}

	signupURL := fmt.Sprintf("%s/signup/admin?token=%s&email=%s",
		strings.TrimSuffix(s.frontendURL, "/"),
		invitation.Token,
		url.QueryEscape(invitation.To),
	)

	data := map[string]interface{}{
		"UserID":    invitation.UserID,
		"Token":     invitation.Token,
		"Email":     invitation.To,
		"ExpiresAt": invitation.ExpiresAt.Format("2006-01-02 15:04:05"),
		"SignupURL": signupURL,
	}

	subject, err := s.loadTemplate("admin_invitation_subject.txt", data)
	if err != nil {
		log.Printf("failed to load admin_invitation_subject template: %v", err)
		subject = "[BookManager] 管理者招待のお知らせ"
	}

	body, err := s.loadTemplate("admin_invitation_body.txt", data)
	if err != nil {
		log.Printf("failed to load admin_invitation_body template: %v", err)
		body = fmt.Sprintf(
			"管理者招待が作成されました。\n\nUserID: %s\nToken: %s\n有効期限: %s\n\n管理者登録画面でトークンを入力してください。\n",
			invitation.UserID,
			invitation.Token,
			invitation.ExpiresAt.Format("2006-01-02 15:04:05"),
		)
	}

	s.sendMail(invitation.To, subject, body)
}

func (s *Service) sendMail(to, subject, body string) {
	if strings.TrimSpace(to) == "" {
		return
	}
	if s.smtp.Host == "" {
		log.Printf("email (dry-run):\nTo: %s\nSubject: %s\n\n%s", to, subject, body)
		return
	}

	from := s.smtp.From
	if strings.TrimSpace(from) == "" {
		from = s.smtp.User
	}
	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s", from, to, subject, body)
	addr := fmt.Sprintf("%s:%s", s.smtp.Host, s.smtp.Port)
	var auth smtp.Auth
	if s.smtp.User != "" && s.smtp.Pass != "" {
		auth = smtp.PlainAuth("", s.smtp.User, s.smtp.Pass, s.smtp.Host)
	}
	if err := smtp.SendMail(addr, auth, from, []string{to}, []byte(msg)); err != nil {
		log.Printf("email send error: %v", err)
	}
}
