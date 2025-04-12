package storage

import (
	"database/sql"
	"fmt"
	"log"
	"rwa/internal/config.go"
)

type Storage struct {
	db *sql.DB
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{
		db: db,
	}
}

func Connect(cfg config.Config) (*sql.DB, error) {
	db, err := sql.Open("sqlite", cfg.StoragePath)
	if err != nil {
		log.Println("failed connected DB")

		return nil, fmt.Errorf("failed connected DB")
	}

	err = db.Ping()
	if err != nil {
		log.Println("failed to connected DB")

		return nil, fmt.Errorf("failed to ping DB")
	}

	return db, nil
}
