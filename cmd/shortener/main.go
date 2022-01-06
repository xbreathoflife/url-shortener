package main

import (
	"github.com/xbreathoflife/url-shortener/internal/app/handler"
	"log"
	"net/http"
)

func main() {
	server := handler.NewURLServer()
	http.HandleFunc("/", server.URLHandler())
	log.Fatal(http.ListenAndServe(":8080", nil))
}
