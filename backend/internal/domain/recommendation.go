package domain

import "time"

type Recommendation struct {
	ID        string    `json:"id"`
	UserID    string    `json:"userId"`
	BookID    string    `json:"bookId"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"createdAt"`
}
