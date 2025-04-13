package app

import (
	"log"
	"net/http"
	"rwa/internal/config.go"
	"rwa/internal/storage/sqlite"
)

// сюда писать код

func GetApp() http.Handler {
	cfg := config.MustLoad()

	log.Println("Config:", cfg.StoragePath, cfg.TokenTTL)

	_, err := sqlite.Connect(cfg)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
