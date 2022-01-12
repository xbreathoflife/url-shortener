package storage

import "fmt"
import "github.com/xbreathoflife/url-shortener/internal/app/entities"

type Storage struct {
	urls   map[int]entities.URL
}

func NewStorage() *Storage {
	storage := &Storage{}
	storage.urls = make(map[int]entities.URL)
	return storage
}

func (storage *Storage) AddURL(baseURL string, shortenedURL string) {
	url := entities.URL{BaseURL: baseURL, ShortenedURL: shortenedURL}

	storage.urls[len(storage.urls)] = url
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
	return len(storage.urls)
}

func (storage *Storage) GetURLIfExist(url string) string {
	for _, v := range storage.urls {
		if v.BaseURL == url {
			return v.ShortenedURL
		}
	}
	return ""
}