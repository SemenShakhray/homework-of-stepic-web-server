package main

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

var (
	ErrUnknownTable   = errors.New("unknown table")
	ErrRecordNotFound = errors.New("record not found")
	ErrInvalidType    = errors.New("have invalid type")
)

func CheckErrors(err error) int {
	if errors.Is(err, ErrUnknownTable) {
		return http.StatusNotFound
	}

	if errors.Is(err, ErrRecordNotFound) {
		return http.StatusNotFound
	}

	if strings.Contains(err.Error(), ErrInvalidType.Error()) {
		return http.StatusBadRequest
	}

	return http.StatusInternalServerError
}

func CreateErrInvalidType(field string) error {
	return fmt.Errorf("field %s have invalid type", field)
}
