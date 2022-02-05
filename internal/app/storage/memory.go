package storage

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
)
import "github.com/xbreathoflife/url-shortener/internal/app/entities"

type MemoryStorage struct {
	urls        map[int]entities.URL
	fileStorage *FileStorage
	BaseURL     string
}

func NewStorage(filePath string, baseURL string) *MemoryStorage {
	storage := &MemoryStorage{}
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

func (storage *MemoryStorage) Init(_ context.Context) error {
	return nil
}

func (storage *MemoryStorage) CheckConnect(_ context.Context) error {
	return nil
}

func (storage *MemoryStorage) InsertNewURL(_ context.Context, id int, baseURL string, shortenedURL string, uuid string) error {
	url := entities.URL{BaseURL: baseURL, ShortenedURL: shortenedURL, UserID: uuid}
	storage.urls[id] = url
	if storage.fileStorage != nil {
		storedURL := entities.StoredURL{
			ID: id,
			BaseURL: baseURL,
			UserID: uuid,
		}
		if err := storage.fileStorage.WriteEvent(storedURL); err != nil {
			return err
		}
	}
	return nil
}

func (storage *MemoryStorage) GetURLByID(_ context.Context, id int) (string, error) {
	url, ok := storage.urls[id]
	if ok {
		return url.BaseURL, nil
	} else {
		return "", fmt.Errorf("URL with id=%d not found", id)
	}
}

func (storage *MemoryStorage) GetUserURLs(_ context.Context, uuid string) ([]entities.URL, error) {
	var urls []entities.URL

	for _, value := range storage.urls {
		if value.UserID == uuid {
			urls = append(urls, value)
		}
	}
	if len(urls) == 0 {
		return nil, errors.New("no urls for this user")
	}
	return urls, nil
}

func (storage *MemoryStorage) GetNextID(_ context.Context) (int, error) {
	return len(storage.urls), nil
}

func (storage *MemoryStorage) GetURLIfExist(_ context.Context, url string) (string, error) {
	for _, v := range storage.urls {
		if v.BaseURL == url {
			return v.ShortenedURL, nil
		}
	}
	return "", nil
}

func (storage *MemoryStorage) GetBaseURL() string {
	return storage.BaseURL
}