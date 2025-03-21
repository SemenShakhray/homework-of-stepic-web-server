package main

import (
	"encoding/json"
	"errors"
	"net/http"
)

var (
	ErrUnknownTable   = errors.New("unknown table")
	ErrRecordNotFound = errors.New("record not found")
)

func CheckErrors(w http.ResponseWriter, resp *Resp) {
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
}
