package core

import (
	"github.com/xbreathoflife/url-shortener/internal/app/storage"
	"strconv"
)

type URLService struct {
	Store *storage.Storage
}

func (us *URLService) GetURLByID(id int) (string, error) {
	return us.Store.GetURL(id)
}

func (us *URLService) AddNewURL(baseURL string) string {
	var shortenedURL string
	if shortenedURL = us.Store.GetURLIfExist(baseURL); shortenedURL == "" {
		urlID := us.Store.GetNextID()
		shortenedURL = us.Store.BaseURL + "/" + strconv.Itoa(urlID)
		us.Store.AddURL(baseURL, shortenedURL)
	}
	return shortenedURL
}
