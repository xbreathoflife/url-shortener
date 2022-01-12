package handler

import (
	"github.com/xbreathoflife/url-shortener/internal/app/core"
	"io"
	"log"
	"net/http"
)

type Handler struct {
	Service *core.URLService
}

func (h *Handler) GetURLHandler(w http.ResponseWriter, r *http.Request, id int) {
	log.Printf("handling get URL at %s\n", r.URL.Path)

	url, err := h.Service.GetURLById(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func(h *Handler) PostURLHandler(w http.ResponseWriter, r *http.Request) {
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

	shortenedURL := h.Service.AddNewURL(baseURL)

	w.WriteHeader(http.StatusCreated)
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	_, err = w.Write([]byte(shortenedURL))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
