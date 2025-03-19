package main

import (
	"database/sql"
	"log"
)

type Storage struct {
	DB *sql.DB
}

func NewStore(db *sql.DB) *Storage {
	return &Storage{
		DB: db,
	}
}
func (s *Storage) GetTables() *Resp {
	stmt, err := s.DB.Prepare("SHOW TABLES")
	if err != nil {
		log.Println(err)

		return &Resp{
			Error: err.Error()}
	}

	rows, err := stmt.Query()
	if err != nil {
		log.Println(err)

		return &Resp{
			Error: err.Error()}
	}
	defer rows.Close()

	var tableName string
	var tableNames []string

	for rows.Next() {
		err := rows.Scan(&tableName)
		if err != nil {
			log.Println(err)

			return &Resp{
				Error: err.Error()}
		}

		tableNames = append(tableNames, tableName)
	}

	if err := rows.Err(); err != nil {
		log.Println(err)

		return &Resp{
			Error: err.Error()}
	}

	return &Resp{
		Response: map[string]interface{}{
			"tables": tableNames,
		},
	}
}
