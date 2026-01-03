package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

var ErrGeminiUnavailable = errors.New("gemini unavailable")

type GeminiClient struct {
	apiKey string
	model  string
	client *http.Client
	prompt string
}

type SeriesGuess struct {
	IsSeries     bool   `json:"isSeries"`
	Name         string `json:"seriesName"`
	VolumeNumber int    `json:"volumeNumber"`
	Confidence   int    `json:"confidence"`
}

type SeriesInput struct {
	Title         string
	Authors       []string
	Publisher     string
	PublishedDate string
	ISBN13        string
	SeriesName    string
}

const defaultPrompt = "You are an assistant that classifies book series from bibliographic data. " +
	"Return JSON only with keys: isSeries (bool), seriesName (string), volumeNumber (number), confidence (0-100). " +
	"If the book is not part of a series, set isSeries=false and seriesName=\"\" and volumeNumber=0."

func NewGeminiClient(apiKey, model, prompt string) *GeminiClient {
	if strings.TrimSpace(prompt) == "" {
		prompt = defaultPrompt
	}
	return &GeminiClient{
		apiKey: apiKey,
		model:  model,
		client: &http.Client{
			Timeout: 6 * time.Second,
		},
		prompt: prompt,
	}
}

func (c *GeminiClient) GuessSeries(ctx context.Context, input SeriesInput) (SeriesGuess, error) {
	if c.apiKey == "" || c.model == "" {
		return SeriesGuess{}, ErrGeminiUnavailable
	}
	endpoint := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s", c.model, c.apiKey)

	prompt := fmt.Sprintf(
		"%s\n\nInput:\n"+
			"title=%q\nauthors=%q\npublisher=%q\npublishedDate=%q\nisbn13=%q\nseriesName=%q\n",
		c.prompt,
		input.Title,
		strings.Join(input.Authors, " / "),
		input.Publisher,
		input.PublishedDate,
		input.ISBN13,
		input.SeriesName,
	)
	payload := map[string]any{
		"contents": []map[string]any{
			{
				"parts": []map[string]string{
					{
						"text": prompt,
					},
				},
			},
		},
		"generationConfig": map[string]any{
			"temperature": 0.1,
		},
	}
	body, _ := json.Marshal(payload)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return SeriesGuess{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return SeriesGuess{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return SeriesGuess{}, ErrGeminiUnavailable
	}
	raw, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return SeriesGuess{}, err
	}

	var data geminiResponse
	if err := json.Unmarshal(raw, &data); err != nil {
		return SeriesGuess{}, err
	}
	text := data.Text()
	var guess SeriesGuess
	if err := json.Unmarshal([]byte(text), &guess); err != nil {
		return SeriesGuess{}, ErrGeminiUnavailable
	}
	if guess.Name != "" && !guess.IsSeries {
		guess.IsSeries = true
	}
	if guess.IsSeries && strings.TrimSpace(guess.Name) == "" {
		return SeriesGuess{}, ErrGeminiUnavailable
	}
	if !guess.IsSeries {
		guess.Name = ""
		guess.VolumeNumber = 0
	}
	return guess, nil
}

type geminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

func (r geminiResponse) Text() string {
	if len(r.Candidates) == 0 {
		return ""
	}
	parts := r.Candidates[0].Content.Parts
	if len(parts) == 0 {
		return ""
	}
	return parts[0].Text
}
