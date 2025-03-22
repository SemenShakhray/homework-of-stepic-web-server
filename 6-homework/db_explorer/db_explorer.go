package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
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
	mux.HandleFunc("/", handler.GetInfoInTable)

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
	var (
		limit, offset int
		err           error
	)

	path := strings.Trim(r.URL.Path, "/")
	args := strings.Split(path, "/")
	if len(args) != 1 || args[0] == "" {
		http.Error(w, ErrUnknownTable.Error(), http.StatusNotFound)
	}
	tableName := args[0]

	exisitsTable := h.store.checkExistsTable(tableName)

	if !exisitsTable {
		http.Error(w, ErrUnknownTable.Error(), http.StatusNotFound)

		return
	}

	offsetString := r.FormValue("offset")
	if offsetString == "" {
		offset = 0
	} else {
		offset, err = strconv.Atoi(offsetString)
		if err != nil {
			offset = 0
		}

	}

	limitString := r.FormValue("limit")
	if limitString == "" {
		limit = 5
	} else {
		limit, err = strconv.Atoi(limitString)
		if err != nil {
			limit = 5
		}

	}

	resp := h.store.GetInfoInTable(tableName, limit, offset)

	CheckErrors(w, resp)

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// тут вы пишете код
// обращаю ваше внимание - в этом задании запрещены глобальные переменные
