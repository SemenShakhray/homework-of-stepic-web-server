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
func (s *Storage) GetListTables() *Resp {
	stmt, err := s.DB.Prepare("SHOW TABLES")
	if err != nil {
		log.Println(err)

		return &Resp{
			Error: err}
	}

	rows, err := stmt.Query()
	if err != nil {
		log.Println(err)

		return &Resp{
			Error: err}
	}
	defer rows.Close()

	var tableName string
	var tableNames []string

	for rows.Next() {
		err := rows.Scan(&tableName)
		if err != nil {
			log.Println(err)

			return &Resp{
				Error: err}
		}

		tableNames = append(tableNames, tableName)
	}

	if err := rows.Err(); err != nil {
		log.Println(err)

		return &Resp{
			Error: err}
	}

	return &Resp{
		Response: map[string]interface{}{
			"tables": tableNames,
		},
	}
}
