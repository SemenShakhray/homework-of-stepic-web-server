package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
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

	return mux, nil

}

func (h *Handler) GetListTables(w http.ResponseWriter, r *http.Request) {
	resp := h.store.GetListTables()

	if resp.Error != nil {
		if errors.Is(resp.Error, ErrUnknownTable) {
			w.WriteHeader(http.StatusNotFound)
			if err := json.NewEncoder(w).Encode(resp); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			return
		}

		if errors.Is(resp.Error, ErrRecordNotFound) {
			w.WriteHeader(http.StatusNotFound)
			if err := json.NewEncoder(w).Encode(resp); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// тут вы пишете код
// обращаю ваше внимание - в этом задании запрещены глобальные переменные
