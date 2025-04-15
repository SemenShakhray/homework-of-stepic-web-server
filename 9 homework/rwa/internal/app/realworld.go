package app

import (
	"log"
	"net/http"

	"rwa/internal/config.go"
	"rwa/internal/handlers"
	"rwa/internal/router"
	"rwa/internal/service"
	"rwa/internal/storage/sqlite"
)

// сюда писать код

func GetApp() http.Handler {
	cfg := config.MustLoad()

	db, err := sqlite.Connect(cfg)
	if err != nil {
		log.Fatal(err)
	}
	// defer func() {
	// 	if errCl := db.Close(); errCl != nil {
	// 		log.Println("Error closing DB: ", errCl)

	// 		return
	// 	}

	// 	log.Println("DB connection closed")
	// }()

	store := sqlite.NewStorage(db)
	service := service.NewService(store)
	handler := handlers.NewHandler(service, cfg)
	router := router.NewRouter(handler)

	return router
}
