package ai

type SeriesGuess struct {
	IsSeries     bool   `json:"isSeries"`
	Name         string `json:"seriesName"`
	VolumeNumber int    `json:"volumeNumber"`
	Confidence   int    `json:"confidence"`
}

type SeriesInput struct {
	Title         string
	RawTitle      string
	Authors       []string
	Publisher     string
	PublishedDate string
	ISBN13        string
	SeriesName    string
}

const defaultPrompt = "You are an assistant that classifies book series from bibliographic data. " +
	"Return JSON only with keys: isSeries (bool), seriesName (string), volumeNumber (number), confidence (0-100). " +
	"If the book is not part of a series, set isSeries=false and seriesName=\"\" and volumeNumber=0."
