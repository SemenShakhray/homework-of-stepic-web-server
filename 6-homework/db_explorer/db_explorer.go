package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
)

type Handler struct {
	store Storage
}

func NewDbExplorer(db *sql.DB) (http.Handler, error) {

	handler := &Handler{
		store: Storage{
			DB: db,
		},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", handler.GetListTables)
	mux.HandleFunc("/tables$", handler.GetInfoInTable)

	return mux, nil

}

func (h *Handler) GetListTables(w http.ResponseWriter, r *http.Request) {
	resp := h.store.GetListTables()

	CheckErrors(w, resp)

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) GetInfoInTable(w http.ResponseWriter, r *http.Request) {
	var limit, offset int
	var err error

	tableName := r.FormValue("tableName")
	if tableName == "" {
		http.Error(w, ErrUnknownTable.Error(), http.StatusNotFound)

		return
	}

	offsetString := r.FormValue("offset")
	if offsetString == "" {
		offset = 0
	} else {
		offset, err = strconv.Atoi(offsetString)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)

			return
		}

	}

	limitString := r.FormValue("limit")
	if limitString == "" {
		limit = 5
	} else {
		limit, err = strconv.Atoi(limitString)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)

			return
		}

	}

	resp := h.store.GetInfoInTable(limit, offset, tableName)

	CheckErrors(w, resp)

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// тут вы пишете код
// обращаю ваше внимание - в этом задании запрещены глобальные переменные
