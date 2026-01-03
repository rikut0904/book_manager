package domain

type ProfileSettings struct {
	UserID        string
	Visibility    string
	GeminiEnabled bool
	GeminiModel   string
	GeminiAPIKey  string
}
