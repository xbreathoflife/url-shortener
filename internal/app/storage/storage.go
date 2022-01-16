package storage

import (
	"fmt"
	"log"
	"strconv"
)
import "github.com/xbreathoflife/url-shortener/internal/app/entities"

type Storage struct {
	urls        map[int]entities.URL
	fileStorage *FileStorage
	BaseURL     string
}

func NewStorage(filePath string, baseURL string) *Storage {
	storage := &Storage{}
	storage.urls = make(map[int]entities.URL)
	storage.BaseURL = baseURL

	var err error
	storage.fileStorage, err = New(filePath)
	if err != nil {
		log.Fatal(err)
	}

	if storage.fileStorage != nil {
		listOfURLs := storage.fileStorage.ReadAllURLsFromFile()
		for i := 0; i < len(listOfURLs); i++ {
			cur := listOfURLs[i]
			storage.urls[cur.ID] = entities.URL{
				BaseURL: cur.BaseURL,
				ShortenedURL: baseURL + "/" + strconv.Itoa(cur.ID),
			}
		}
	}

	return storage
}

func (storage *Storage) AddURL(baseURL string, shortenedURL string) {
	url := entities.URL{BaseURL: baseURL, ShortenedURL: shortenedURL}
	id := len(storage.urls)
	storage.urls[id] = url
	if storage.fileStorage != nil {
		storedURL := entities.StoredURL{
			ID: id,
			BaseURL: baseURL,
		}
		if err := storage.fileStorage.WriteEvent(storedURL); err != nil {
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