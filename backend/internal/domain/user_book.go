package domain

type UserBook struct {
	ID           string `json:"id"`
	UserID       string `json:"userId"`
	BookID       string `json:"bookId"`
	Note         string `json:"note"`
	AcquiredAt   string `json:"acquiredAt"`
	SeriesID     string `json:"seriesId"`
	VolumeNumber int    `json:"volumeNumber"`
}
