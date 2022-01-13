package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/xbreathoflife/url-shortener/internal/app/core"
	"github.com/xbreathoflife/url-shortener/internal/app/handler"
	"github.com/xbreathoflife/url-shortener/internal/app/storage"
	"net/http"
)

type urlServer struct {
	store *storage.Storage
	handlers *handler.Handler
}

func NewURLServer() *urlServer {
	store := storage.NewStorage()
	handlers := handler.Handler{Service: &core.URLService{Store: store} }
	return &urlServer{store: store, handlers: &handlers}
}

func (us *urlServer) URLHandler() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/", func(rw http.ResponseWriter, r *http.Request) {
		us.handlers.PostURLHandler(rw, r)
	})
	r.Get("/{urlID}", func(rw http.ResponseWriter, r *http.Request) {
		urlID := chi.URLParam(r, "urlID")
		us.handlers.GetURLHandler(rw, r, urlID)
	})

	r.MethodNotAllowed(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "Wrong path", http.StatusBadRequest)
	})

	return r
}