package sqlite

import (
	"database/sql"
	"fmt"
	"log"

	"rwa/internal/config.go"
	"rwa/internal/storage"

	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func NewStorage(db *sql.DB) storage.Storer {
	return &Storage{
		db: db,
	}
}

func Connect(cfg config.Config) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", cfg.StoragePath)
	if err != nil {
		log.Println("failed connected DB:", err)

		return nil, fmt.Errorf("failed connected DB: %w", err)
	}

	err = db.Ping()
	if err != nil {
		log.Println("failed to connected DB")

		return nil, fmt.Errorf("failed to ping DB")
	}

	return db, nil
}
