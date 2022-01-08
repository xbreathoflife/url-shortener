package handler

import (
	"github.com/xbreathoflife/url-shortener/internal/app/storage"
	"io"
	"log"
	"net/http"
	"strconv"
)

const serverHost = "http://localhost:8080"

func GetURLHandler(w http.ResponseWriter, r *http.Request, id int, store *storage.Storage) {
	log.Printf("handling get URL at %s\n", r.URL.Path)

	url, err := store.GetURL(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func PostURLHandler(w http.ResponseWriter, r *http.Request, store *storage.Storage) {
	log.Printf("handling post URL at %s\n", r.URL.Path)

	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	baseURL := string(b)
	if baseURL == "" {
		http.Error(w, "Empty body - no url", http.StatusBadRequest)
		return
	}

	var shortenedURL string
	if shortenedURL = store.GetURLIfExist(baseURL); shortenedURL == "" {
		urlID := store.GetNextID()
		shortenedURL = serverHost + "/" + strconv.Itoa(urlID)
		store.AddURL(baseURL, shortenedURL)
	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	_, err = w.Write([]byte(shortenedURL))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
