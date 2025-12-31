package domain

import "time"

type Recommendation struct {
	ID        string
	UserID    string
	BookID    string
	Comment   string
	CreatedAt time.Time
}
