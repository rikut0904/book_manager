package domain

type Book struct {
	ID            string   `json:"id"`
	ISBN13        string   `json:"isbn13"`
	Title         string   `json:"title"`
	OriginalTitle string   `json:"originalTitle"`
	Authors       []string `json:"authors"`
	Publisher     string   `json:"publisher"`
	PublishedDate string   `json:"publishedDate"`
	ThumbnailURL  string   `json:"thumbnailUrl"`
	Source        string   `json:"source"`
	SeriesName    string   `json:"seriesName"`
}
