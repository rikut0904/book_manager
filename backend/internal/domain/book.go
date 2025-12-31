package domain

type Book struct {
	ID            string
	ISBN13        string
	Title         string
	Authors       []string
	Publisher     string
	PublishedDate string
	ThumbnailURL  string
	Source        string
}
