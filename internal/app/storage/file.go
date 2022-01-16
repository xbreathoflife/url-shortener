package storage

import (
	"encoding/json"
	"github.com/xbreathoflife/url-shortener/internal/app/entities"
	"io"
	"log"
	"os"
)

type FileStorage struct {
	file    *os.File
	encoder *json.Encoder
	decoder *json.Decoder
}

func New(filePath string) (*FileStorage, error) {
	if filePath == "" {
		return nil, nil
	}

	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return nil, err
	}

	return &FileStorage{
		file:    file,
		encoder: json.NewEncoder(file),
		decoder: json.NewDecoder(file),
	}, nil
}

func (s *FileStorage) WriteEvent(u entities.URL) error {
	return s.encoder.Encode(u)
}

func (s *FileStorage) ReadAllURLsFromFile() []entities.URL {
	var listOfURLs []entities.URL

	for {
		var u entities.URL
		err := s.decoder.Decode(&u)
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		listOfURLs = append(listOfURLs, u)
	}
	return listOfURLs
}

func (s *FileStorage) Close() error {
	return s.file.Close()
}
