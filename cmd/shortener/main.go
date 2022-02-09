package main

import (
	"flag"
	"github.com/xbreathoflife/url-shortener/config"
	"github.com/xbreathoflife/url-shortener/internal/app/server"
	"github.com/xbreathoflife/url-shortener/internal/app/storage"
	"log"
	"net/http"
)

func parseFlags(conf *config.Config) {
	address := flag.String("a", "", "Адрес запуска HTTP-сервера")
	baseURL := flag.String("b", "", "Базовый адрес результирующего сокращённого URL")
	filePath := flag.String("f", "", "Путь до файла с сокращёнными URL")
	connString := flag.String("d", "", "Строка с адресом подключения к БД")
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

	if *connString != "" {
		conf.ConnString = *connString
	}
}

func main() {
	conf := config.Init()
	parseFlags(&conf)
	dbStorage := storage.NewDBStorage(conf.ConnString, conf.BaseURL)
	urlServer := server.NewURLServer(dbStorage)
	r := urlServer.URLHandler()
	log.Fatal(http.ListenAndServe(conf.Address, r))
}
