package domain

type BookTag struct {
	UserID string `json:"userId"`
	BookID string `json:"bookId"`
	TagID  string `json:"tagId"`
}
