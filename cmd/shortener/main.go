package main

import (
	"github.com/xbreathoflife/url-shortener/internal/app/server"
	"log"
	"net/http"
)

func main() {
	urlServer := server.NewURLServer()
	r := urlServer.URLHandler()
	log.Fatal(http.ListenAndServe(":8080", r))
}
