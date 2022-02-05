package entities

type URL struct {
	BaseURL      string `json:"original_url"`
	ShortenedURL string `json:"short_url"`
	UserID       string `json:"-"`
}

type StoredURL struct {
	UserID  string `json:"uuid"`
	BaseURL string `json:"url"`
	ID      int    `json:"id"`
}

type BaseURL struct {
	Name string `json:"url"`
}

type ShortenedURL struct {
	Name string `json:"result"`
}
