package handler

import (
	"github.com/xbreathoflife/url-shortener/internal/app/storage"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

const serverHost = "http://localhost:8080"

type urlServer struct {
	store *storage.Storage
}

func NewURLServer() *urlServer {
	store := storage.NewStorage()
	return &urlServer{store: store}
}

func (us *urlServer) getURLHandler(w http.ResponseWriter, r *http.Request, id int) {
	log.Printf("handling get URL at %s\n", r.URL.Path)

	url, err := us.store.GetURL(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (us *urlServer) postURLHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("handling post URL at %s\n", r.URL.Path)

	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	baseURL := string(b)
	var shortenedURL string
	if shortenedURL = us.store.GetURLIfExist(baseURL); shortenedURL == "" {
		urlID := us.store.GetNextID()
		shortenedURL = serverHost + "/" + strconv.Itoa(urlID)
		us.store.AddURL(baseURL, shortenedURL)
	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	_, err = w.Write([]byte(shortenedURL))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (us *urlServer) URLHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			if r.Method != http.MethodPost {
				http.Error(w, "Only POST requests are allowed!", http.StatusBadRequest)
				return
			}
			us.postURLHandler(w, r)
		} else {
			if r.Method != http.MethodGet {
				http.Error(w, "Only GET requests are allowed!", http.StatusBadRequest)
				return
			}
			path := strings.Trim(r.URL.Path, "/")
			id, err := strconv.Atoi(path)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			us.getURLHandler(w, r, id)
		}
	}
}

