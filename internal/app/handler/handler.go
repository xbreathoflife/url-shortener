package handler

import (
	"encoding/json"
	"github.com/xbreathoflife/url-shortener/internal/app/auth"
	"github.com/xbreathoflife/url-shortener/internal/app/core"
	"github.com/xbreathoflife/url-shortener/internal/app/entities"
	"io"
	"log"
	"net/http"
	"strconv"
)

type Handler struct {
	Service *core.URLService
}

func (h *Handler) GetURLHandler(w http.ResponseWriter, r *http.Request, urlID string) {
	log.Printf("handling get URL at %s\n", r.URL.Path)

	if urlID == "" {
		http.Error(w, "urlID param is missed", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(urlID)
	if err != nil {
		http.Error(w, "urlID must be an integer", http.StatusBadRequest)
		return
	}

	url, err := h.Service.GetURLByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *Handler) GetUserURLs(w http.ResponseWriter, r *http.Request) {
	log.Printf("handling get user URLs at %s\n", r.URL.Path)

	uuid := r.Context().Value(auth.CtxKey).(string)

	URLsForUser, err := h.Service.GetUserURLs(uuid)
	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	js, err := json.Marshal(URLsForUser)
	if err != nil {
		http.Error(w, "Error during building response json", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(js)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func(h *Handler) PostURLHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("handling post URL at %s\n", r.URL.Path)
	uuid := r.Context().Value(auth.CtxKey).(string)

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

	shortenedURL := h.Service.AddNewURL(baseURL, uuid)

	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(shortenedURL))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func(h *Handler) PostJSONURLHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("handling post URL at %s\n", r.URL.Path)

	uuid := r.Context().Value(auth.CtxKey).(string)
	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	baseURL := entities.BaseURL{}
	if err := json.Unmarshal(b, &baseURL); err != nil {
		http.Error(w, "Error during parsing request json", http.StatusBadRequest)
		return
	}
	if baseURL.Name == "" {
		http.Error(w, "Empty body - no url", http.StatusBadRequest)
		return
	}

	shortURL := h.Service.AddNewURL(baseURL.Name, uuid)
	shortenedURL := entities.ShortenedURL{Name: shortURL}
	js, err := json.Marshal(shortenedURL)
	if err != nil {
		http.Error(w, "Error during building response json", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(js)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}