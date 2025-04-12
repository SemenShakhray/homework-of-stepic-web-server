package app

import (
	"log"
	"net/http"
	"rwa/internal/config.go"
	"rwa/internal/storage"
)

// сюда писать код

func GetApp() http.Handler {
	cfg := config.MustLoad()

	db, err := storage.Connect(cfg)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
