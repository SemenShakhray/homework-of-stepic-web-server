package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
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
	mux.HandleFunc("/", handler.DinamicServe)

	return mux, nil

}

func (h *Handler) DinamicServe(w http.ResponseWriter, r *http.Request) {

	path := strings.Trim(r.URL.Path, "/")
	if path == "" {
		h.GetListTables(w, r)

		return
	}

	arr := strings.Split(path, "/")
	switch r.Method {
	case http.MethodGet:
		if len(arr) == 1 {
			h.GetInfoInTable(w, r, arr)
		}
		if len(arr) == 2 {
			h.GetInfoRecord(w, r, arr)
		}
	}
}

func (h *Handler) GetListTables(w http.ResponseWriter, r *http.Request) {
	resp, err := h.store.GetListTables()

	if err != nil {
		code := CheckErrors(err)

		responseError(w, err, code)

		return
	}

	responseOK(w, resp)
}

func (h *Handler) GetInfoInTable(w http.ResponseWriter, r *http.Request, arr []string) {
	var (
		limit, offset = 5, 0
	)

	if arr[0] == "" {
		responseError(w, ErrUnknownTable, http.StatusNotFound)

		return
	}
	tableName := arr[0]

	exisitsTable := h.store.checkExistsTable(tableName)

	if !exisitsTable {
		responseError(w, ErrUnknownTable, http.StatusNotFound)

		return
	}

	offsetString := r.FormValue("offset")
	if offsetString != "" {
		offset, _ = strconv.Atoi(offsetString)
	}

	limitString := r.FormValue("limit")
	if limitString != "" {
		limit, _ = strconv.Atoi(limitString)
	}

	resp, err := h.store.GetInfoInTable(tableName, limit, offset)

	if err != nil {
		code := CheckErrors(err)
		responseError(w, err, code)

		return
	}

	responseOK(w, resp)
}

func (h *Handler) GetInfoRecord(w http.ResponseWriter, r *http.Request, arr []string) {
	if arr[0] == "" {
		responseError(w, ErrUnknownTable, http.StatusNotFound)

		return
	}
	tableName := arr[0]

	exisitsTable := h.store.checkExistsTable(tableName)

	if !exisitsTable {
		responseError(w, ErrUnknownTable, http.StatusNotFound)

		return
	}

	if arr[1] == "" {
		responseError(w, ErrRecordNotFound, http.StatusNotFound)

		return
	}
	id, err := strconv.Atoi(arr[1])
	if err != nil {
		responseError(w, fmt.Errorf("failed id"), http.StatusBadRequest)
	}

	resp, err := h.store.GetInfoRecord(tableName, id)
	if err != nil {
		code := CheckErrors(err)

		responseError(w, err, code)

		return
	}

	responseOK(w, resp)
}

func responseError(w http.ResponseWriter, err error, code int) {
	resp := Resp{
		Error: err.Error(),
	}

	w.WriteHeader(code)
	Err := json.NewEncoder(w).Encode(resp)
	if Err != nil {
		http.Error(w, Err.Error(), http.StatusInternalServerError)
	}
}

func responseOK(w http.ResponseWriter, resp interface{}) {
	w.WriteHeader(http.StatusOK)
	Err := json.NewEncoder(w).Encode(resp)
	if Err != nil {
		http.Error(w, Err.Error(), http.StatusInternalServerError)
	}
}

// тут вы пишете код
// обращаю ваше внимание - в этом задании запрещены глобальные переменные
