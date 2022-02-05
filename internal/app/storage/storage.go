package storage

import (
	"errors"
	"fmt"
	"log"
	"strconv"
)
import "github.com/xbreathoflife/url-shortener/internal/app/entities"

type UserStorage map[int]entities.URL

type Storage struct {
	urls        map[string]UserStorage
	fileStorage *FileStorage
	BaseURL     string
}

func NewStorage(filePath string, baseURL string) *Storage {
	storage := &Storage{}
	storage.urls = make(map[string]UserStorage)
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
			if len(storage.urls[cur.UserID]) == 0 {
				storage.urls[cur.UserID] = make(UserStorage)
			}
			storage.urls[cur.UserID][cur.ID] = entities.URL{
				BaseURL: cur.BaseURL,
				ShortenedURL: baseURL + "/" + strconv.Itoa(cur.ID),
			}
		}
	}

	return storage
}

func (storage *Storage) AddURL(baseURL string, shortenedURL string, uuid string) {
	url := entities.URL{BaseURL: baseURL, ShortenedURL: shortenedURL}
	id := len(storage.urls[uuid])
	if id == 0 {
		storage.urls[uuid] = make(UserStorage)
	}
	storage.urls[uuid][id] = url
	if storage.fileStorage != nil {
		storedURL := entities.StoredURL{
			ID: id,
			BaseURL: baseURL,
			UserID: uuid,
		}
		if err := storage.fileStorage.WriteEvent(storedURL); err != nil {
			log.Fatal(err)
		}
	}
}

func (storage *Storage) GetURL(id int, uuid string) (string, error) {
	url, ok := storage.urls[uuid][id]
	if ok {
		return url.BaseURL, nil
	} else {
		return "", fmt.Errorf("URL with id=%d not found", id)
	}
}

func (storage *Storage) GetUserURLs(uuid string) (UserStorage, error) {
	urls, ok := storage.urls[uuid]
	if ok {
		return urls, nil
	} else {
		return nil, errors.New("no urls for this user")
	}
}

func (storage *Storage) GetNextID(uuid string) int {
	return len(storage.urls[uuid])
}

func (storage *Storage) GetURLIfExist(url string, uuid string) string {
	for _, v := range storage.urls[uuid] {
		if v.BaseURL == url {
			return v.ShortenedURL
		}
	}
	return ""
}