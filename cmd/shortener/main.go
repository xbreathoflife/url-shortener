package main

import (
	"flag"
	"github.com/xbreathoflife/url-shortener/config"
	"github.com/xbreathoflife/url-shortener/internal/app/server"
	"log"
	"net/http"
)

func parseFlags(conf *config.Config) {
	address := flag.String("a", "", "Адрес запуска HTTP-сервера")
	baseURL := flag.String("b", "", "Базовый адрес результирующего сокращённого URL")
	filePath := flag.String("f", "", "Путь до файла с сокращёнными URL")
	flag.Parse()

	if *address != "" {
		conf.Address = *address
	}

	if *baseURL != "" {
		conf.BaseURL = *baseURL
	}

	if *filePath != "" {
		conf.FilePath = *filePath
	}
}

func main() {
	conf := config.Init()
	parseFlags(&conf)
	urlServer := server.NewURLServer(conf.BaseURL, conf.FilePath)
	r := urlServer.URLHandler()
	log.Fatal(http.ListenAndServe(conf.Address, r))
}
