package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
)

type infoTable struct {
	Field   string
	Type    string
	Null    string
	Default *string
	Extra   string
}

type Storage struct {
	DB *sql.DB
}

func NewStore(db *sql.DB) *Storage {
	return &Storage{
		DB: db,
	}
}
func (s *Storage) GetListTables() (*Resp, error) {
	stmt, err := s.DB.Prepare("SHOW TABLES")
	if err != nil {
		log.Println(err)

		return nil, err

	}

	rows, err := stmt.Query()
	if err != nil {
		log.Println(err)

		return nil, err
	}
	defer rows.Close()

	var tableName string
	var tableNames []string

	for rows.Next() {
		err := rows.Scan(&tableName)
		if err != nil {
			log.Println(err)

			return nil, err
		}

		tableNames = append(tableNames, tableName)
	}

	if err := rows.Err(); err != nil {
		log.Println(err)

		return nil, err
	}

	return &Resp{
		Response: map[string]interface{}{
			"tables": tableNames,
		},
	}, nil
}

func (s *Storage) GetInfoInTable(tableName string, limit, offset int) (*Resp, error) {
	query := fmt.Sprintf("SELECT * FROM %s LIMIT ? OFFSET ?", tableName)

	rows, err := s.DB.Query(query, limit, offset)
	if err != nil {
		if errors.Is(err, fmt.Errorf("%s doesn't exist", tableName)) {
			log.Println(err)

			return nil, ErrUnknownTable
		}

		return nil, err
	}
	defer rows.Close()

	namesColumns, err := rows.Columns()
	if err != nil {
		log.Println("Error when getting column names:", err)
		return nil, err
	}

	countColumns := len(namesColumns)

	result := make([]map[string]interface{}, 0)

	tempLineInterface := make([]interface{}, countColumns)

	pTempLineInterface := make([]interface{}, countColumns)
	for i := 0; i < countColumns; i++ {
		pTempLineInterface[i] = &tempLineInterface[i]
	}

	tempLine := make([]interface{}, countColumns)

	// Обрабатываем строки данных
	for rows.Next() {
		resMap := make(map[string]interface{}, countColumns)

		err := rows.Scan(pTempLineInterface...)
		if err != nil {
			log.Println("failed scanning string:", err)
			return nil, err
		}

		for i, value := range tempLineInterface {
			switch value.(type) {
			case []byte:
				tempLine[i] = string(value.([]byte))
			default:
				tempLine[i] = *(pTempLineInterface[i].(*interface{}))
			}
			resMap[namesColumns[i]] = tempLine[i]
		}
		result = append(result, resMap)
	}

	if err := rows.Err(); err != nil {
		log.Println("failed scanning string:", err)

		return nil, err
	}

	return &Resp{
		Response: map[string]interface{}{
			"records": result,
		},
	}, nil

}

func (s *Storage) GetInfoRecord(tableName string, id int) (*Resp, error) {
	query := fmt.Sprintf("SELECT * FROM %s LIMIT 1", tableName)

	rows, err := s.DB.Query(query)
	if err != nil {
		if errors.Is(err, fmt.Errorf("%s doesn't exist", tableName)) {
			log.Println(err)

			return nil, ErrUnknownTable
		}
		return nil, err
	}

	namesColumns, err := rows.Columns()
	if err != nil {
		log.Println("Error when getting column names:", err)
		return nil, err
	}
	rows.Close()

	countColumns := len(namesColumns)

	resMap := make(map[string]interface{})
	tempLine := make([]interface{}, countColumns)
	pTempLine := make([]interface{}, countColumns)

	for i, _ := range tempLine {
		pTempLine[i] = &tempLine[i]
	}

	query = fmt.Sprintf("SELECT * FROM %s where id = ?", tableName)
	row := s.DB.QueryRow(query, id)

	err = row.Scan(pTempLine...)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("record don't exists")

			return nil, ErrRecordNotFound
		} else {
			log.Println("failed scanning record")

			return nil, err
		}
	}

	for i, value := range tempLine {
		switch value.(type) {
		case []byte:
			resMap[namesColumns[i]] = string(value.([]byte))
		default:
			resMap[namesColumns[i]] = *(pTempLine[i].(*interface{}))
		}
	}

	return &Resp{
		Response: map[string]interface{}{
			"record": resMap,
		},
	}, nil
}

func (s *Storage) AddItem(data map[string]interface{}, tableName string) (*Resp, error) {
	// var infoFull []infoTable

	// query := fmt.Sprintf("SHOW FULL COLUMNS FROM %s", tableName)
	// rows, err := s.DB.Query(query)
	// if err != nil {
	// 	return nil, err
	// }
	// defer rows.Close()

	// for rows.Next() {
	// 	var info infoTable
	// 	err := rows.Scan(&info.Field, &info.Type, new(interface{}), &info.Null, new(interface{}), &info.Default, &info.Extra, new(interface{}), new(interface{}))
	// 	if err != nil {
	// 		log.Println("failed scanning string:", err)
	// 		return nil, err
	// 	}
	// 	infoFull = append(infoFull, info)
	// }

	// if err := rows.Err(); err != nil {
	// 	log.Println("failed scanning string:", err)
	// 	return nil, err
	// }

	// log.Println(infoFull)

	var keys, placeHolders []string
	var values []interface{}

	for key, value := range data {
		if key == "id" {
			continue
		} else {
			keys = append(keys, key)
			values = append(values, value)
		}
	}

	for i := 0; i < len(keys); i++ {
		placeHolders = append(placeHolders, "?")
	}
	log.Println(tableName, keys, placeHolders, values)
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", tableName, strings.Join(keys, ","), strings.Join(placeHolders, ","))

	log.Println(query)

	row, err := s.DB.Exec(query, values...)
	if err != nil {
		log.Println("failed add record: ", err)

		return nil, err
	}

	id, err := row.LastInsertId()
	if err != nil {
		log.Println("don't received id: ", err)

		return nil, err
	}

	return &Resp{
		Response: map[string]interface{}{
			"id": id,
		},
	}, nil
}

func (s *Storage) UpdateRecord(tableName string, data map[string]any, id int) (*Resp, error) {
	var keys []string
	var values []any

	for key, value := range data {
		log.Printf("key: %s, value: %v, type: %T", key, value, value)
		if key == "id" {
			return nil, fmt.Errorf("field id have invalid type")
		} else {
			keys = append(keys, key+"=?")
			values = append(values, value)
		}
	}
	values = append(values, id)

	query := fmt.Sprintf("UPDATE %s SET %s WHERE id = ?", tableName, strings.Join(keys, ","))
	log.Println(query)

	row, err := s.DB.Exec(query, values...)
	if err != nil {
		log.Println("failed update record: ", err)

		return nil, err
	}

	count, err := row.RowsAffected()
	if err != nil {
		if err != nil {
			log.Println("no updates: ", err)

			return nil, err
		}
	}

	return &Resp{
		Response: map[string]any{
			"updated": count,
		},
	}, nil
}

func (s *Storage) checkExistsTable(tableName string) bool {

	query := fmt.Sprintf("SHOW TABLES LIKE '%s'", tableName)
	rows, err := s.DB.Query(query)
	if err != nil {
		log.Println(err)
		return false
	}
	defer rows.Close()

	return rows.Next()
}
