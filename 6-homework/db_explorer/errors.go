package main

import "errors"

var (
	ErrUnknownTable   = errors.New("unknown table")
	ErrRecordNotFound = errors.New("record not found")
)
