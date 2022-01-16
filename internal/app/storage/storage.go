package storage

import (
	"fmt"
	"log"
)
import "github.com/xbreathoflife/url-shortener/internal/app/entities"

type Storage struct {
	urls   map[int]entities.URL
	fileStorage *FileStorage
}

func NewStorage(filePath string) *Storage {
	storage := &Storage{}
	storage.urls = make(map[int]entities.URL)

	var err error
	storage.fileStorage, err = New(filePath)
	if err != nil {
		log.Fatal(err)
	}

	if storage.fileStorage != nil {
		listOfURLs := storage.fileStorage.ReadAllURLsFromFile()
		for i := 0; i < len(listOfURLs); i++ {
			storage.urls[i] = listOfURLs[i]
		}
	}

	return storage
}

func (storage *Storage) AddURL(baseURL string, shortenedURL string) {
	url := entities.URL{BaseURL: baseURL, ShortenedURL: shortenedURL}

	storage.urls[len(storage.urls)] = url
	if storage.fileStorage != nil {
		if err := storage.fileStorage.WriteEvent(url); err != nil {
			log.Fatal(err)
		}
	}
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