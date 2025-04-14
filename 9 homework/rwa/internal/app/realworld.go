package app

import (
	"log"
	"net/http"
	"rwa/internal/config.go"
	"rwa/internal/handlers"
	"rwa/internal/storage/sqlite"
	"rwa/router"
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
	handler := handlers.NewHandler(store)
	router := router.NewRouter(handler)

	return router
}
