package worker

import (
	"context"
	"github.com/xbreathoflife/url-shortener/internal/app/entities"
	"github.com/xbreathoflife/url-shortener/internal/app/storage"
	"log"
	"time"
)
const bufferSize = 10

type DeleteWorker struct {
	deleteBuffer chan entities.DeleteTask
	ctx          context.Context
}

func NewDeleteWorker(ctx context.Context) *DeleteWorker {
	return &DeleteWorker{deleteBuffer: make(chan entities.DeleteTask, bufferSize), ctx: ctx}
}

func (dw *DeleteWorker) AddUrlForDeleting(task entities.DeleteTask) {
	dw.deleteBuffer <- task
}

func (dw *DeleteWorker) RunDeleting(storage storage.Storage) {
	ticker := time.NewTicker(3 * time.Second)
	log.Println("here")
	for {
		select {
		case <-ticker.C:
			var items []entities.DeleteTask
			for i := 0; i < bufferSize; i++ {
				select {
				case item := <-dw.deleteBuffer:
					log.Println(item.Uuid, item.ShortURLID, "added")
					items = append(items, item)
				default:
				}
			}
			if len(items) > 0 {
				log.Println("deleting")
				err := storage.DeleteBatch(dw.ctx, items)
				if err != nil {
					log.Println(err)
				}
			}
		case <-dw.ctx.Done():
			return
		}
	}
}



