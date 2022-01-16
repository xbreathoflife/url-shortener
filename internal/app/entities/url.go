package entities

type URL struct {
	BaseURL      string `json:"base"`
	ShortenedURL string `json:"shortened"`
}

type BaseURL struct {
	Name string `json:"url"`
}

type ShortenedURL struct {
	Name string `json:"result"`
}
