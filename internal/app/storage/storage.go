package storage

import "fmt"

type URL struct {
	BaseURL      string
	ShortenedURL string
}

type Storage struct {
	urls   map[int]URL
	nextID int
}

func NewStorage() *Storage {
	storage := &Storage{}
	storage.urls = make(map[int]URL)
	storage.nextID = 0
	return storage
}

func (storage *Storage) AddURL(baseURL string, shortenedURL string) {
	url := URL{baseURL, shortenedURL}

	storage.urls[storage.nextID] = url
	storage.nextID++
}

func (storage *Storage) GetURL(id int) (string, error) {
	url, ok := storage.urls[id]
	if ok {
		return url.BaseURL, nil
	} else {
		return "", fmt.Errorf("URL with id=%d not found", id)
	}
}

func (storage *Storage) GetNextID() int {
	return storage.nextID
}

func (storage *Storage) GetURLIfExist(url string) string {
	for _, v := range storage.urls {
		if v.BaseURL == url {
			return v.ShortenedURL
		}
	}
	return ""
}