package domain

type Favorite struct {
	ID       string `json:"id"`
	UserID   string `json:"userId"`
	Type     string `json:"type"`
	BookID   string `json:"bookId"`
	SeriesID string `json:"seriesId"`
}
