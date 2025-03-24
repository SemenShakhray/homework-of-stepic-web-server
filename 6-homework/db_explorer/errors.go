package main

import (
	"errors"
	"net/http"
)

var (
	ErrUnknownTable   = errors.New("unknown table")
	ErrRecordNotFound = errors.New("record not found")
)

func CheckErrors(err error) int {
	if errors.Is(err, ErrUnknownTable) {
		return http.StatusNotFound
	}

	if errors.Is(err, ErrRecordNotFound) {
		return http.StatusNotFound
	}

	return http.StatusInternalServerError
}
