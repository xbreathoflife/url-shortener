package entities

type URL struct {
	BaseURL      string
	ShortenedURL string
}

type BaseURL struct {
	Name string `json:"url"`
}

type ShortenedURL struct {
	Name string `json:"result"`
}
