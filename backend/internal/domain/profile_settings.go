package domain

type ProfileSettings struct {
	UserID        string
	Visibility    string
	OpenAIEnabled bool
	OpenAIModel   string
	OpenAIAPIKey  string
}
