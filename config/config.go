package config

import (
	"github.com/caarlos0/env/v6"
	"log"
)

type Config struct {
	Address string `env:"SERVER_ADDRESS"`
	BaseURL string `env:"BASE_URL"`
}

func Init() Config {
	cfg := Config{
		Address: ":8080",
		BaseURL: "http://localhost:8080",
	}
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	return cfg
}