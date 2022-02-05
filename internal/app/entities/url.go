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

type Record struct {
	ID           int    `json:"-"`
	BaseURL      string `json:"-"`
	ShortenedURL string `json:"short_url"`
	UserID       string `json:"-"`
	CorID        string `json:"correlation_id"`
}

type BatchURLRequest struct {
	Name string `json:"original_url"`
	ID   string `json:"correlation_id"`
}

type ShortenedURL struct {
	Name string `json:"result"`
}
