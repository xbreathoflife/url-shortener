package handler

import (
	"encoding/json"
	"errors"
	"github.com/xbreathoflife/url-shortener/internal/app/auth"
	"github.com/xbreathoflife/url-shortener/internal/app/core"
	"github.com/xbreathoflife/url-shortener/internal/app/entities"
	er "github.com/xbreathoflife/url-shortener/internal/app/errors"
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

	ctx := r.Context()
	url, err := h.Service.GetURLByID(ctx, id)
	if err != nil {
		var nfe *er.ULRNotFoundError
		if errors.As(err, &nfe) {
			w.WriteHeader(http.StatusGone)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *Handler) GetUserURLs(w http.ResponseWriter, r *http.Request) {
	log.Printf("handling get user URLs at %s\n", r.URL.Path)

	ctx := r.Context()
	uuid := ctx.Value(auth.CtxKey).(string)

	URLsForUser, err := h.Service.GetUserURLs(ctx, uuid)
	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		var se *er.EmptyStorageError
		if errors.As(err, &se) {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	js, err := json.Marshal(URLsForUser)
	if err != nil {
		http.Error(w, "Error during building response json", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(js)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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

	ctx := r.Context()
	uuid := ctx.Value(auth.CtxKey).(string)
	shortenedURL, err := h.Service.AddNewURL(ctx, baseURL, uuid)
	var statusCode int
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")

	if err != nil {
		var de *er.ULRDuplicateError
		if errors.As(err, &de) {
			shortenedURL = de.ShortURL
			statusCode = http.StatusConflict
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		statusCode = http.StatusCreated
	}

	w.WriteHeader(statusCode)
	_, err = w.Write([]byte(shortenedURL))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func(h *Handler) PostJSONURLHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("handling post URL at %s\n", r.URL.Path)

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

	ctx := r.Context()
	uuid := ctx.Value(auth.CtxKey).(string)
	var statusCode int
	shortURL, err := h.Service.AddNewURL(ctx, baseURL.Name, uuid)
	if err != nil {
		var de *er.ULRDuplicateError
		if errors.As(err, &de) {
			shortURL = de.ShortURL
			statusCode = http.StatusConflict
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		statusCode = http.StatusCreated
	}

	shortenedURL := entities.ShortenedURL{Name: shortURL}
	js, err := json.Marshal(shortenedURL)
	if err != nil {
		http.Error(w, "Error during building response json", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, err = w.Write(js)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func(h *Handler) PostJSONURLBatchHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("handling post URL at %s\n", r.URL.Path)

	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var urls []entities.BatchURLRequest
	if err := json.Unmarshal(b, &urls); err != nil {
		http.Error(w, "Error during parsing request json", http.StatusBadRequest)
		return
	}

	if len(urls) == 0 {
		http.Error(w, "Empty body", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	uuid := ctx.Value(auth.CtxKey).(string)
	records, err := h.Service.AddURLsBatch(ctx, urls, uuid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	js, err := json.Marshal(records)
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

func (h *Handler) GetPing(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	err := h.Service.Storage.CheckConnect(ctx)
	if err != nil {
		http.Error(w, "Couldn't connect to DB", http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) DeleteURLs(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var ids []string
	if err := json.Unmarshal(b, &ids); err != nil {
		http.Error(w, "Error during parsing request json", http.StatusBadRequest)
		return
	}
	for _, i := range ids {
		log.Printf(i)
	}

	ctx := r.Context()
	uuid := ctx.Value(auth.CtxKey).(string)

	h.Service.AsyncDelete(ctx, uuid, ids)
	w.WriteHeader(http.StatusAccepted)
}