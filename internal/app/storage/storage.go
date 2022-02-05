package storage

import (
	"context"
	"github.com/xbreathoflife/url-shortener/internal/app/entities"
)

type Storage interface {
	Init(ctx context.Context) error
	CheckConnect(ctx context.Context) error
	InsertNewURL(ctx context.Context, id int, baseURL string, shortenedURL string, uuid string) error
	GetURLByID(ctx context.Context, id int) (string, error)
	GetUserURLs(ctx context.Context, uuid string) ([]entities.URL, error)
	GetNextID(ctx context.Context) (int, error)
	GetURLIfExist(ctx context.Context, url string) (string, error)
	GetBaseURL() string
}
