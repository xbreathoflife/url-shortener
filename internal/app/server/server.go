package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/xbreathoflife/url-shortener/internal/app/auth"
	"github.com/xbreathoflife/url-shortener/internal/app/compress"
	"github.com/xbreathoflife/url-shortener/internal/app/core"
	"github.com/xbreathoflife/url-shortener/internal/app/handler"
	"github.com/xbreathoflife/url-shortener/internal/app/storage"
	"net/http"
)

type urlServer struct {
	store     *storage.Storage
	dbStorage *storage.DBStorage
	handlers  *handler.Handler
}

func NewURLServer(baseURL string, filePath string, connString string) *urlServer {
	store := storage.NewStorage(filePath, baseURL)
	dbStorage := storage.NewDBStorage(connString)
	handlers := handler.Handler{Service: &core.URLService{Store: store, DBStore: dbStorage}}
	return &urlServer{store: store, dbStorage: dbStorage, handlers: &handlers}
}

func (us *urlServer) URLHandler() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(compress.GzipDecoder)
	r.Use(compress.GzipEncoder)
	r.Use(auth.AuthToken)

	r.Post("/", func(rw http.ResponseWriter, r *http.Request) {
		us.handlers.PostURLHandler(rw, r)
	})

	r.Get("/{urlID}", func(rw http.ResponseWriter, r *http.Request) {
		urlID := chi.URLParam(r, "urlID")
		us.handlers.GetURLHandler(rw, r, urlID)
	})

	r.Post("/api/shorten", func(rw http.ResponseWriter, r *http.Request) {
		us.handlers.PostJSONURLHandler(rw, r)
	})

	r.Get("/user/urls", func(rw http.ResponseWriter, r *http.Request) {
		us.handlers.GetUserURLs(rw, r)
	})

	r.Get("/ping", func(rw http.ResponseWriter, r *http.Request) {
		us.handlers.GetPing(rw, r)
	})

	r.MethodNotAllowed(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "Wrong path", http.StatusBadRequest)
	})

	return r
}