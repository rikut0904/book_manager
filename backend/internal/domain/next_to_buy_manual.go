package domain

type NextToBuyManual struct {
	ID           string `json:"id"`
	UserID       string `json:"userId"`
	Title        string `json:"title"`
	SeriesName   string `json:"seriesName"`
	VolumeNumber int    `json:"volumeNumber"`
	Note         string `json:"note"`
}
