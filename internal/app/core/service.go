package core

import (
	"context"
	"github.com/xbreathoflife/url-shortener/internal/app/entities"
	"github.com/xbreathoflife/url-shortener/internal/app/storage"
	"strconv"
)

type URLService struct {
	Storage   storage.Storage
}

func (us *URLService) GetURLByID(ctx context.Context, id int) (string, error) {
	return us.Storage.GetURLByID(ctx, id)
}

func (us *URLService) GetUserURLs(ctx context.Context, uuid string) ([]entities.URL, error) {
	return us.Storage.GetUserURLs(ctx, uuid)
}

func (us *URLService) AddNewURL(ctx context.Context, baseURL string, uuid string) (string, error) {
	shortenedURL, err := us.Storage.GetURLIfExist(ctx, baseURL)
	if err != nil {
		return "", err
	}
	if shortenedURL == "" {
		urlID, err := us.Storage.GetNextID(ctx)
		if err != nil {
			return "", err
		}
		shortenedURL = us.Storage.GetBaseURL() + "/" + strconv.Itoa(urlID)
		err = us.Storage.InsertNewURL(ctx, urlID, baseURL, shortenedURL, uuid)
		if err != nil {
			return "", err
		}
	}

	return shortenedURL, nil
}
