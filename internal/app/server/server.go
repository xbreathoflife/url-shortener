package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/xbreathoflife/url-shortener/internal/app/handler"
	"github.com/xbreathoflife/url-shortener/internal/app/storage"
	"net/http"
	"strconv"
)

type urlServer struct {
	store *storage.Storage
}

func NewURLServer() *urlServer {
	store := storage.NewStorage()
	return &urlServer{store: store}
}

func (us *urlServer) URLHandler() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/", func(rw http.ResponseWriter, r *http.Request) {
		handler.PostURLHandler(rw, r, us.store)
	})
	r.Get("/{urlID}", func(rw http.ResponseWriter, r *http.Request) {
		urlID := chi.URLParam(r, "urlID")
		if urlID == "" {
			http.Error(rw, "urlID param is missed", http.StatusBadRequest)
			return
		}
		id, err := strconv.Atoi(urlID)
		if err != nil {
			http.Error(rw, "urlID must be an integer", http.StatusBadRequest)
			return
		}
		handler.GetURLHandler(rw, r, id, us.store)
	})

	r.MethodNotAllowed(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "Wrong path", http.StatusBadRequest)
	})

	return r
}