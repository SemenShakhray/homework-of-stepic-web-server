package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"reflect"
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
		log.Println("failed getting column names:", err)
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

	idString, err := s.NameFieldWithID(tableName)
	if err != nil {
		return nil, err
	}

	query = fmt.Sprintf("SELECT * FROM %s where %s = ?", tableName, idString)
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
	err := s.checkTypeParamRequest(tableName, data)
	if err != nil {
		return nil, err
	}

	var keys, placeHolders []string
	var values []interface{}

	idString, err := s.NameFieldWithID(tableName)
	if err != nil {
		return nil, err
	}

	for key, value := range data {
		if key == idString {
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
	err := s.checkTypeParamRequest(tableName, data)
	if err != nil {
		return nil, err
	}

	var keys []string
	var values []any

	idString, err := s.NameFieldWithID(tableName)
	if err != nil {
		return nil, err
	}

	for key, value := range data {
		log.Printf("key: %s, value: %v, type: %T", key, value, value)
		if key == idString {
			return nil, fmt.Errorf("field %s have invalid type", idString)
		} else {
			keys = append(keys, key+"=?")
			values = append(values, value)
		}
	}
	values = append(values, id)

	query := fmt.Sprintf("UPDATE %s SET %s WHERE %s = ?", tableName, strings.Join(keys, ","), idString)
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

func (s *Storage) DeleteRecord(tableName string, id int) (*Resp, error) {
	idString, err := s.NameFieldWithID(tableName)
	if err != nil {
		return nil, err
	}

	query := fmt.Sprintf("DELETE FROM %s WHERE %s = ?", tableName, idString)
	row, err := s.DB.Exec(query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}

	count, err := row.RowsAffected()
	if err != nil {
		return nil, err
	}

	return &Resp{
		Response: map[string]interface{}{
			"deleted": count,
		},
	}, nil
}

func (s *Storage) checkTypeParamRequest(tableName string, data map[string]interface{}) error {
	sqlToGoTypeMap := map[string]reflect.Kind{
		"VARCHAR": reflect.String,
		"TEXT":    reflect.String,
		"CHAR":    reflect.String,
		"INT":     reflect.Int,
		"FLOAT":   reflect.Float64,
	}

	query := fmt.Sprintf("SELECT * FROM %s LIMIT 1", tableName)
	rows, err := s.DB.Query(query)
	if err != nil {
		log.Println("failed request DB: ", err)

		return err
	}
	defer rows.Close()

	columnInfo, err := rows.ColumnTypes()
	if err != nil {
		log.Println("failed reseived info of columns: ", err)

		return err
	}

	for _, col := range columnInfo {
		colType := col.DatabaseTypeName()

		if v, ok := data[col.Name()]; ok {
			if data[col.Name()] == nil {
				nullable, _ := col.Nullable()
				if !nullable {
					return fmt.Errorf("field %s have invalid type", col.Name())
				}
				continue
			}

			switch reflect.TypeOf(v).Kind() {
			case reflect.Float64:
				v = int(v.(float64))
			}
			log.Println("colNAme:", col.Name(), "typeV:", reflect.TypeOf(v), "colType:", colType)
			if reflect.TypeOf(v).Kind() != sqlToGoTypeMap[colType] {
				return fmt.Errorf("field %s have invalid type", col.Name())
			}
		}
	}

	return nil
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

func (s *Storage) NameFieldWithID(tableName string) (string, error) {
	query := fmt.Sprintf("SELECT * FROM %s LIMIT 1", tableName)
	rows, err := s.DB.Query(query)
	if err != nil {
		log.Println("failed request DB: ", err)

		return "", err
	}
	defer rows.Close()

	columnName, err := rows.Columns()
	if err != nil {
		log.Println("failed reseived info of columns: ", err)

		return "", err
	}

	return columnName[0], nil
}
