package main

import (
	"database/sql"
	"errors"
	"fmt"
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

func (s *Storage) GetInfoInTable(tableName string, offset, limit int) *Resp {

	stmt, err := s.DB.Prepare("SELECT * FROM" + tableName + "LIMIT ? OFFSET ?")
	if err != nil {
		log.Println(err)

		return &Resp{
			Error: err}
	}

	rows, err := stmt.Query(limit, offset)
	if err != nil {
		if errors.Is(err, fmt.Errorf("%s doesn't exist", tableName)) {
			log.Println(err)

			return &Resp{
				Error: ErrUnknownTable}
		}

		return &Resp{
			Error: err,
		}
	}
	defer rows.Close()

	var columns []string
	columns, err = rows.Columns()
	if err != nil {
		log.Println(err)

		return &Resp{
			Error: err}
	}

	var values []interface{}
	values = make([]interface{}, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			log.Println(err)

			return &Resp{
				Error: err}
		}

		for i, col := range values {
			if col != nil {
				if value, ok := col.([]byte); ok {
					values[i] = string(value)
				} else {
					values[i] = col
				}
			}
		}
	}

	return &Resp{
		Response: map[string]interface{}{
			"columns": columns,
			"values":  values,
		},
	}

}

func (s *Storage) checkExistsTable(tableName string) bool {
	stmt, err := s.DB.Prepare("SHOW TABLES LIKE ?")
	if err != nil {
		log.Println(err)

		return false
	}

	rows, err := stmt.Query(tableName)
	if err != nil {
		log.Println(err)
		return false
	}
	defer rows.Close()

	for rows.Next() {
		return true
	}

	return false
}
