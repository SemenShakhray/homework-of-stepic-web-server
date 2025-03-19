package main

import (
	"database/sql"
	"net/http"
)

// type Storage struct {
// 	db *sql.DB
// }

func NewDbExplorer(db *sql.DB) (http.Handler, error) {

	store := &Storage{db: db}

	mux := http.NewServeMux()
	mux.HandleFunc("/", store.GetTables)

	return nil, nil

}

// func (s *Storage) GetTables(w http.ResponseWriter, r *http.Request) {
// 	stmt, err := s.db.Prepare("SHOW TABLES")
// 	if err != nil {
// 		log.Println(err)

// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}

// 	type CR map[string]interface{}

// 	rows, err := stmt.Query()
// 	if err != nil {
// 		log.Println(err)

// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}
// 	defer rows.Close()

// 	for rows.Next() {

// 	}
// 	err := rows.Scan()

// }

// тут вы пишете код
// обращаю ваше внимание - в этом задании запрещены глобальные переменные
