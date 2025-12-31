package domain

type UserBook struct {
	ID           string
	UserID       string
	BookID       string
	Note         string
	AcquiredAt   string
	SeriesID     string
	VolumeNumber int
}
