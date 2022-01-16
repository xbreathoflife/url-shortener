package entities

type URL struct {
	BaseURL      string
	ShortenedURL string
}

type StoredURL struct {
	BaseURL string `json:"url"`
	ID      int    `json:"id"`
}

type BaseURL struct {
	Name string `json:"url"`
}

type ShortenedURL struct {
	Name string `json:"result"`
}
