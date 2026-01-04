package domain

import "time"

type OpenAIKey struct {
	ID        string
	Name      string
	APIKey    string
	CreatedAt time.Time
}
