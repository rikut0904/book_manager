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

var ErrOpenAIUnavailable = errors.New("openai unavailable")

type OpenAIClient struct {
	apiKey string
	model  string
	client *http.Client
	prompt string
}

func NewOpenAIClient(apiKey, model, prompt string) *OpenAIClient {
	if strings.TrimSpace(prompt) == "" {
		prompt = defaultPrompt
	}
	model = strings.TrimSpace(model)
	if model == "" {
		model = "gpt-4o-mini"
	}
	return &OpenAIClient{
		apiKey: apiKey,
		model:  model,
		client: &http.Client{
			Timeout: 8 * time.Second,
		},
		prompt: prompt,
	}
}

type openAIResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func (c *OpenAIClient) GuessSeries(ctx context.Context, input SeriesInput) (SeriesGuess, error) {
	if c.apiKey == "" || c.model == "" {
		return SeriesGuess{}, ErrOpenAIUnavailable
	}
	payload := map[string]any{
		"model": c.model,
		"messages": []map[string]string{
			{
				"role":    "system",
				"content": c.prompt,
			},
			{
				"role": "user",
				"content": fmt.Sprintf(
					"Input:\n"+
						"title=%q\nrawTitle=%q\nauthors=%q\npublisher=%q\npublishedDate=%q\nisbn13=%q\nseriesName=%q\n",
					input.Title,
					input.RawTitle,
					strings.Join(input.Authors, " / "),
					input.Publisher,
					input.PublishedDate,
					input.ISBN13,
					input.SeriesName,
				),
			},
		},
		"response_format": map[string]string{
			"type": "json_object",
		},
		"temperature": 0.1,
	}
	body, _ := json.Marshal(payload)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.openai.com/v1/chat/completions", bytes.NewReader(body))
	if err != nil {
		return SeriesGuess{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.client.Do(req)
	if err != nil {
		return SeriesGuess{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		raw, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<16))
		return SeriesGuess{}, fmt.Errorf("openai status %d: %s", resp.StatusCode, strings.TrimSpace(string(raw)))
	}
	raw, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return SeriesGuess{}, err
	}
	var data openAIResponse
	if err := json.Unmarshal(raw, &data); err != nil {
		return SeriesGuess{}, err
	}
	if len(data.Choices) == 0 {
		return SeriesGuess{}, ErrOpenAIUnavailable
	}
	text := data.Choices[0].Message.Content
	var guess SeriesGuess
	if err := json.Unmarshal([]byte(text), &guess); err != nil {
		return SeriesGuess{}, fmt.Errorf("openai parse error: %w text=%q", err, text)
	}
	if guess.Name != "" && !guess.IsSeries {
		guess.IsSeries = true
	}
	if guess.IsSeries && strings.TrimSpace(guess.Name) == "" {
		return SeriesGuess{}, ErrOpenAIUnavailable
	}
	if !guess.IsSeries {
		guess.Name = ""
		guess.VolumeNumber = 0
	}
	return guess, nil
}
