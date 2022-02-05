package core

import (
	"github.com/xbreathoflife/url-shortener/internal/app/entities"
	"github.com/xbreathoflife/url-shortener/internal/app/storage"
	"strconv"
)

type URLService struct {
	Store *storage.Storage
}

func (us *URLService) GetURLByID(id int) (string, error) {
	return us.Store.GetURL(id)
}

func (us *URLService) GetUserURLs(uuid string) ([]entities.URL, error) {
	return us.Store.GetUserURLs(uuid)
}

func (us *URLService) AddNewURL(baseURL string, uuid string) string {
	var shortenedURL string
	if shortenedURL = us.Store.GetURLIfExist(baseURL); shortenedURL == "" {
		urlID := us.Store.GetNextID()
		shortenedURL = us.Store.BaseURL + "/" + strconv.Itoa(urlID)
		us.Store.AddURL(baseURL, shortenedURL, uuid)
	}
	return shortenedURL
}
