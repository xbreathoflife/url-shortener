package core

import (
	"context"
	"fmt"
	"github.com/xbreathoflife/url-shortener/internal/app/entities"
	"github.com/xbreathoflife/url-shortener/internal/app/storage"
	"github.com/xbreathoflife/url-shortener/internal/app/worker"
	"strconv"
)

type URLService struct {
	Storage   storage.Storage
	DeleteWorker *worker.DeleteWorker
}

func (us *URLService) GetURLByID(ctx context.Context, id int) (string, error) {
	return us.Storage.GetURLByID(ctx, id)
}

func (us *URLService) GetUserURLs(ctx context.Context, uuid string) ([]entities.URL, error) {
	return us.Storage.GetUserURLs(ctx, uuid)
}

func (us *URLService) AddNewURL(ctx context.Context, baseURL string, uuid string) (string, error) {
	urlID, err := us.Storage.GetNextID(ctx)
	if err != nil {
		return "", err
	}

	shortenedURL := us.Storage.GetBaseURL() + "/" + strconv.Itoa(urlID)
	err = us.Storage.InsertNewURL(ctx, urlID, baseURL, shortenedURL, uuid)

	if err != nil {
		return "", err
	}

	return shortenedURL, nil
}

func (us *URLService) AddURLsBatch(ctx context.Context, urls []entities.BatchURLRequest, uuid string) ([]entities.Record, error) {
	urlID, err := us.Storage.GetNextID(ctx)
	if err != nil {
		return nil, err
	}
	records := make([]entities.Record, 0, len(urls))
	for i, u := range urls {
		var record entities.Record
		record.ID = i + urlID
		record.BaseURL = u.Name
		record.CorID = u.ID
		record.UserID = uuid
		record.ShortenedURL = us.Storage.GetBaseURL() + "/" + strconv.Itoa(record.ID)
		records = append(records, record)
	}
	err = us.Storage.InsertBatch(ctx, records)
	if err != nil {
		return nil, err
	}

	return records, nil
}

func (us *URLService) AsyncDelete(_ context.Context, uuid string, ids []string) {
	go func() {
		for _, id := range ids {
			newID, err := strconv.Atoi(id)
			if err != nil {
				fmt.Printf("error: %v\n", err)
				return
			}
			us.DeleteWorker.AddURLForDeleting(entities.DeleteTask{ShortURLID: newID, UUID: uuid})
		}
	}()
}