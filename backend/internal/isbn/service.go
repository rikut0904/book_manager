package isbn

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"book_manager/backend/internal/domain"
	"book_manager/backend/internal/repository"
)

var ErrNotFound = errors.New("book not found")

type Service struct {
	client  *http.Client
	baseURL string
	apiKey  string
	ttl     time.Duration
	cache   repository.IsbnCacheRepository
}

type SeriesGuess struct {
	Name         string
	VolumeNumber int
}

func NewService(baseURL, apiKey string, ttl time.Duration, cache repository.IsbnCacheRepository) *Service {
	return &Service{
		client: &http.Client{
			Timeout: 8 * time.Second,
		},
		baseURL: baseURL,
		apiKey:  apiKey,
		ttl:     ttl,
		cache:   cache,
	}
}

func (s *Service) Lookup(isbn string) (domain.Book, SeriesGuess, error) {
	if book, ok := s.fromCache(isbn); ok {
		series := inferSeries(book.Title, book.SeriesName)
		return book, series, nil
	}

	book, err := s.fetchGoogleBooks(isbn)
	if err != nil {
		return domain.Book{}, SeriesGuess{}, err
	}
	s.storeCache(isbn, book)
	series := inferSeries(book.Title, book.SeriesName)
	return book, series, nil
}

func (s *Service) fromCache(isbn string) (domain.Book, bool) {
	entry, ok := s.cache.Get(isbn)
	if !ok {
		return domain.Book{}, false
	}
	if s.ttl > 0 && time.Since(entry.FetchedAt) > s.ttl {
		return domain.Book{}, false
	}
	return entry.Book, true
}

func (s *Service) storeCache(isbn string, book domain.Book) {
	_ = s.cache.Upsert(domain.IsbnCache{
		ISBN13:    isbn,
		Book:      book,
		FetchedAt: time.Now(),
	})
}

func (s *Service) fetchGoogleBooks(isbn string) (domain.Book, error) {
	query := url.Values{}
	query.Set("q", fmt.Sprintf("isbn:%s", isbn))
	if s.apiKey != "" {
		query.Set("key", s.apiKey)
	}

	requestURL := fmt.Sprintf("%s?%s", s.baseURL, query.Encode())
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		return domain.Book{}, err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return domain.Book{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return domain.Book{}, fmt.Errorf("google books status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return domain.Book{}, err
	}

	var payload googleBooksResponse
	if err := json.Unmarshal(body, &payload); err != nil {
		return domain.Book{}, err
	}
	if len(payload.Items) == 0 {
		return domain.Book{}, ErrNotFound
	}

	item := payload.Items[0]
	book := domain.Book{
		ID:            item.ID,
		ISBN13:        extractISBN13(item.VolumeInfo.IndustryIdentifiers, isbn),
		Title:         item.VolumeInfo.Title,
		Authors:       item.VolumeInfo.Authors,
		Publisher:     item.VolumeInfo.Publisher,
		PublishedDate: item.VolumeInfo.PublishedDate,
		ThumbnailURL:  item.VolumeInfo.ImageLinks.Thumbnail,
		Source:        "google",
		SeriesName:    item.VolumeInfo.Series,
	}
	return book, nil
}

func extractISBN13(identifiers []industryIdentifier, fallback string) string {
	for _, id := range identifiers {
		if id.Type == "ISBN_13" && id.Identifier != "" {
			return id.Identifier
		}
	}
	return fallback
}

type googleBooksResponse struct {
	Items []googleBooksItem `json:"items"`
}

type googleBooksItem struct {
	ID         string            `json:"id"`
	VolumeInfo googleBooksVolume `json:"volumeInfo"`
}

type googleBooksVolume struct {
	Title               string                `json:"title"`
	Authors             []string              `json:"authors"`
	Publisher           string                `json:"publisher"`
	PublishedDate       string                `json:"publishedDate"`
	IndustryIdentifiers []industryIdentifier  `json:"industryIdentifiers"`
	ImageLinks          googleBooksImageLinks `json:"imageLinks"`
	Series              string                `json:"series"`
}

type industryIdentifier struct {
	Type       string `json:"type"`
	Identifier string `json:"identifier"`
}

type googleBooksImageLinks struct {
	Thumbnail string `json:"thumbnail"`
}

var (
	volumePattern         = regexp.MustCompile(`(?i)(?:第?\s*([0-9０-９]+)\s*(?:巻|冊|話)|vol\.?\s*([0-9０-９]+))`)
	trailingNumberPattern = regexp.MustCompile(`\s*([0-9０-９]+)\s*$`)
	bracketPattern = regexp.MustCompile(`【[^】]+】`)
	parenPattern = regexp.MustCompile(`[（(][^)）]+[)）]`)
)

func inferSeries(title, seriesName string) SeriesGuess {
	guess := SeriesGuess{
		Name:         strings.TrimSpace(seriesName),
		VolumeNumber: 0,
	}
	if title == "" {
		return guess
	}
	if guess.Name == "" {
		guess.Name = strings.TrimSpace(volumePattern.ReplaceAllString(title, ""))
		guess.Name = strings.TrimSpace(trailingNumberPattern.ReplaceAllString(guess.Name, ""))
		guess.Name = strings.Trim(guess.Name, " -‐–—・")
	}
	return guess
}

func InferSeries(title, seriesName string) SeriesGuess {
	return inferSeries(title, seriesName)
}

func NormalizeSeriesName(name string) string {
	cleaned := strings.TrimSpace(name)
	if cleaned == "" {
		return ""
	}
	cleaned = bracketPattern.ReplaceAllString(cleaned, "")
	cleaned = parenPattern.ReplaceAllString(cleaned, "")
	cleaned = volumePattern.ReplaceAllString(cleaned, "")
	cleaned = trailingNumberPattern.ReplaceAllString(cleaned, "")
	cleaned = strings.TrimSpace(cleaned)
	cleaned = strings.Trim(cleaned, " -‐–—・")
	return strings.TrimSpace(cleaned)
}

func NormalizeTitle(title string) string {
	cleaned := strings.TrimSpace(title)
	if cleaned == "" {
		return ""
	}
	cleaned = bracketPattern.ReplaceAllString(cleaned, "")
	cleaned = parenPattern.ReplaceAllString(cleaned, "")
	cleaned = volumePattern.ReplaceAllString(cleaned, "")
	cleaned = trailingNumberPattern.ReplaceAllString(cleaned, "")
	cleaned = strings.TrimSpace(cleaned)
	cleaned = strings.Trim(cleaned, " -‐–—・")
	return strings.TrimSpace(cleaned)
}
