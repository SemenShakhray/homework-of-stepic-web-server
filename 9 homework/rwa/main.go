package main

import (
	"fmt"
	"net/http"
	"rwa/internal/app"
)

// сюда код писать не надо

func main() {
	addr := ":8080"
	h := app.GetApp()
	fmt.Println("start server at", addr)
	http.ListenAndServe(addr, h)
}
