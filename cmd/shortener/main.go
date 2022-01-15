package main

import (
	"github.com/xbreathoflife/url-shortener/config"
	"github.com/xbreathoflife/url-shortener/internal/app/server"
	"log"
	"net/http"
)

func main() {
	conf := config.Init()
	urlServer := server.NewURLServer(conf.BaseURL)
	r := urlServer.URLHandler()
	log.Fatal(http.ListenAndServe(conf.Address, r))
}
