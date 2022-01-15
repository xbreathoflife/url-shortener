package entities

type URL struct {
	BaseURL      string
	ShortenedURL string
}

type BaseURL struct {
	Url string `json:"url"`
}

type ShortenedURL struct {
	Url string `json:"result"`
}
