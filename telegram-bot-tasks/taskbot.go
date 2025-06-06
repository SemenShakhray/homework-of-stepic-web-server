package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"taskbot/internal/config"
	"taskbot/internal/handlers"
	"taskbot/internal/storage/sqlite"
	"taskbot/service"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func startTaskBot(ctx context.Context, httpListenAddr string) error {
	conf := config.MustLoad()

	db, err := sqlite.Connect(conf)
	if err != nil {
		return err
	}

	store := sqlite.NewStorage(db)
	service := service.NewService(store)
	handler := handlers.NewHandler(service)
	_, err = handlers.NewBot(handler, conf)
	if err != nil {
		return err
	}

	fmt.Println("start listen", httpListenAddr)
	return http.ListenAndServe(httpListenAddr, nil)
}

func main() {
	err := startTaskBot(context.Background(), ":8081")
	if err != nil {
		log.Fatalln(err)
	}
}

// это заглушка чтобы импорт сохранился
func __dummy() {
	tgbotapi.APIEndpoint = "_dummy"
}
