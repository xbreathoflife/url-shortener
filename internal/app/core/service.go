package core

import (
	"github.com/xbreathoflife/url-shortener/internal/app/storage"
	"strconv"
)

type URLService struct {
	Store *storage.Storage
}

func (us *URLService) GetURLByID(id int, uuid string) (string, error) {
	return us.Store.GetURL(id, uuid)
}

func (us *URLService) GetUserURLs(uuid string) (storage.UserStorage, error) {
	return us.Store.GetUserURLs(uuid)
}

func (us *URLService) AddNewURL(baseURL string, uuid string) string {
	var shortenedURL string
	if shortenedURL = us.Store.GetURLIfExist(baseURL, uuid); shortenedURL == "" {
		urlID := us.Store.GetNextID(uuid)
		shortenedURL = us.Store.BaseURL + "/" + strconv.Itoa(urlID)
		us.Store.AddURL(baseURL, shortenedURL, uuid)
	}
	return shortenedURL
}
