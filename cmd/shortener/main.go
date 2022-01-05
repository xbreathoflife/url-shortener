package main

import (
	"github.com/xbreathoflife/url-shortener/internal/app"
	"log"
	"net/http"
)

func main() {
	server := app.NewURLServer()
	http.HandleFunc("/", server.TaskHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
