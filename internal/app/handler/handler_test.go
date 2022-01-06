package handler

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestURLPostHandler(t *testing.T) {
	type want struct {
		contentType string
		statusCode  int
		url         []string
	}
	tests := []struct {
		name    string
		request string
		body    []string
		want    want
	}{
		{
			name: "multiple posts",
			body: []string{"https://yandex.ru/", "https://www.google.ru/", "https://www.youtube.com/"},
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  201,
				url:         []string{"http://localhost:8080/0", "http://localhost:8080/1", "http://localhost:8080/2"},
			},
			request: "/",
		},
		{
			name: "repeated url",
			body: []string{"https://yandex.ru/", "https://www.google.ru/", "https://yandex.ru/"},
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  201,
				url:         []string{"http://localhost:8080/0", "http://localhost:8080/1", "http://localhost:8080/0"},
			},
			request: "/",
		},
		{
			name: "wrong path #1",
			body: []string{"https://yandex.ru/"},
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  400,
				url:         []string{"Only GET requests are allowed!\n"},
			},
			request: "/0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := NewURLServer()
			for i, element := range tt.body {
				body := []byte(element)
				request := httptest.NewRequest(http.MethodPost, tt.request, bytes.NewBuffer(body))
				w := httptest.NewRecorder()
				h := server.URLHandler()
				h.ServeHTTP(w, request)
				result := w.Result()

				assert.Equal(t, tt.want.statusCode, result.StatusCode)
				assert.Equal(t, tt.want.contentType, w.Header().Get("Content-Type"))

				urlResult, err := ioutil.ReadAll(result.Body)
				require.NoError(t, err)
				err = result.Body.Close()
				require.NoError(t, err)

				assert.Equal(t, tt.want.url[i], string(urlResult))
			}
		})
	}
}

func TestURLGetHandler(t *testing.T) {
	type want struct {
		statusCode int
		url        []string
	}
	tests := []struct {
		name    string
		request string
		body    []string
		want    want
	}{
		{
			name: "multiple post and get",
			body: []string{"https://yandex.ru/", "https://www.google.ru/", "https://www.youtube.com/"},
			want: want{
				statusCode: 307,
				url:        []string{"http://localhost:8080/0", "http://localhost:8080/1", "http://localhost:8080/2"},
			},
			request: "/",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := NewURLServer()
			for i, element := range tt.body {
				body := []byte(element)
				request := httptest.NewRequest(http.MethodPost, tt.request, bytes.NewBuffer(body))
				w := httptest.NewRecorder()
				h := server.URLHandler()
				h.ServeHTTP(w, request)
				target := tt.request + strconv.Itoa(i)
				getRequest := httptest.NewRequest(http.MethodGet, target, nil)
				w = httptest.NewRecorder()
				h.ServeHTTP(w, getRequest)
				result := w.Result()

				assert.Equal(t, tt.want.statusCode, result.StatusCode)
				assert.Equal(t, element, w.Header().Get("Location"))

				err := result.Body.Close()
				require.NoError(t, err)
			}
		})
	}
}

func TestURLHandlerError(t *testing.T) {
	tests := []struct {
		name       string
		request    string
		body       string
		method     string
		statusCode int
	}{
		{
			name:       "get instead of post",
			body:       "https://yandex.ru/",
			statusCode: 400,
			method:     http.MethodGet,
			request:    "/",
		},
		{
			name:       "not number in get",
			body:       "https://yandex.ru/",
			statusCode: 400,
			method:     http.MethodGet,
			request:    "/notnumber",
		},
		{
			name:       "post instead of get",
			body:       "https://yandex.ru/",
			statusCode: 400,
			method:     http.MethodPost,
			request:    "/1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := NewURLServer()
			body := []byte(tt.body)
			request := httptest.NewRequest(tt.method, tt.request, bytes.NewBuffer(body))
			w := httptest.NewRecorder()
			h := server.URLHandler()
			h.ServeHTTP(w, request)
			result := w.Result()
			assert.Equal(t, tt.statusCode, result.StatusCode)
			err := result.Body.Close()
			require.NoError(t, err)
		})
	}
}
