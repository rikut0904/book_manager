package domain

import "time"

type IsbnCache struct {
	ISBN13    string
	Book      Book
	FetchedAt time.Time
}
