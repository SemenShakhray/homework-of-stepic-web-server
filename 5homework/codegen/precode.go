package main

// import (
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// 	"strconv"
// )

// type Result map[string]interface{}

// func (h *MyApi) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// 	switch r.URL.Path {
// 	case "/user/profile":
// 		h.GetUser(w, r)
// 	case "/user/create":
// 		h.CreateUser(w, r)
// 	default:
// 		http.Error(w, `{"error":"unknown method"}`, 404)
// 	}
// }

// func (h *MyApi) GetUser(w http.ResponseWriter, r *http.Request) {
// 	if r.Method == http.MethodPost || r.Method == http.MethodGet {
// 		login := r.FormValue("login")
// 		if login == "" {
// 			http.Error(w, `{"error": "login must me not empty"}`, 400)
// 			return
// 		}
// 		in := ProfileParams{
// 			Login: login,
// 		}
// 		ctx := context.TODO()
// 		res, err := h.Profile(ctx, in)
// 		if err != nil {
// 			ae, ok := err.(ApiError)
// 			if !ok {
// 				http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
// 				return
// 			}
// 			http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), ae.HTTPStatus)
// 			return
// 		}
// 		result := Result{
// 			"error": "",
// 			"response": Result{
// 				"id":        res.ID,
// 				"login":     res.Login,
// 				"full_name": res.FullName,
// 				"status":    res.Status,
// 			},
// 		}
// 		err = json.NewEncoder(w).Encode(result)
// 		if err != nil {
// 			http.Error(w, "error encoding JSON", http.StatusInternalServerError)
// 		}
// 	}
// }

// func (h *MyApi) CreateUser(w http.ResponseWriter, r *http.Request) {
// 	if r.Method == http.MethodPost {
// 		password := r.Header.Get("X-Auth")
// 		if password == "" {
// 			http.Error(w, `{"error":"unauthorized"}`, 403)
// 			return
// 		}
// 		if password != "100500" {
// 			http.Error(w, `{"error":"wrong password"}`, http.StatusUnauthorized)
// 			return
// 		}
// 		in := CreateParams{}
// 		in.Login = r.FormValue("login")
// 		if in.Login == "" {
// 			http.Error(w, `{"error":"login must me not empty"}`, 400)
// 			return
// 		}
// 		if in.Login != "" && len(in.Login) < 10 {
// 			http.Error(w, `{"error":"login len must be >= 10"}`, 400)
// 			return
// 		}
// 		in.Status = r.FormValue("status")
// 		if in.Status == "" {
// 			in.Status = "user"
// 		}

// 		if !(in.Status == "user" || in.Status == "moderator" || in.Status == "admin") {
// 			http.Error(w, `{"error":"status must be one of [user, moderator, admin]"}`, 400)
// 			return
// 		}

// 		in.Name = r.FormValue("full_name")

// 		age, err := strconv.Atoi(r.FormValue("age"))
// 		if err != nil {
// 			http.Error(w, `{"error":"age must be int"}`, 400)
// 			return
// 		}
// 		in.Age = age
// 		if in.Age < 0 {
// 			http.Error(w, `{"error":"age must be >= 0"}`, 400)
// 			return
// 		}
// 		if in.Age > 128 {
// 			http.Error(w, `{"error":"age must be <= 128"}`, 400)
// 			return
// 		}
// 		ctx := context.TODO()
// 		res, err := h.Create(ctx, in)
// 		if err != nil {
// 			ae, ok := err.(ApiError)
// 			if !ok {
// 				http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
// 				return
// 			}
// 			http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), ae.HTTPStatus)
// 			return
// 		}
// 		result := Result{
// 			"error": "",
// 			"response": Result{
// 				"id": res.ID,
// 			}}
// 		err = json.NewEncoder(w).Encode(result)
// 		if err != nil {
// 			http.Error(w, "error encoding JSON", http.StatusInternalServerError)
// 		}
// 	} else {
// 		http.Error(w, `{"error":"bad method"}`, 406)
// 	}
// }
