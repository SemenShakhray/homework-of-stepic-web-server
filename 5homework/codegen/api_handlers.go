package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type Result map[string]interface{}

func (h *MyApi) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/user/profile":
		h.ProfileUser(w, r)
	case "/user/create":
		h.CreateUser(w, r)
	default:
		http.Error(w, "{\"error\":\"unknown method\"}", 404)
	}
}
func (h *OtherApi) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/user/create":
		h.CreateUser(w, r)
	default:
		http.Error(w, "{\"error\":\"unknown method\"}", 404)
	}
}

func ValidProfileParams(r *http.Request) (ProfileParams, error) {
	in := ProfileParams{}

	in.Login = r.FormValue("login")
	if in.Login == "" {
		return in, fmt.Errorf("{\"error\":\"login must me not empty\"}")
	}
	return in, nil
}

func ValidCreateParams(r *http.Request) (CreateParams, error) {
	in := CreateParams{}

	in.Login = r.FormValue("login")
	if in.Login == "" {
		return in, fmt.Errorf("{\"error\":\"login must me not empty\"}")
	}
	if len(in.Login) < 10 {
		return in, fmt.Errorf("{\"error\":\"login len must be >= 10\"}")
	}

	in.Name = r.FormValue("full_name")

	in.Status = r.FormValue("status")
	if in.Status == "" {
		in.Status = "user"
	}
	enum := strings.Split("user|moderator|admin", "|")
	mapEnum := make(map[string]string)
	for _, v := range enum {
		mapEnum[v] = ""
	}
	if _, ok := mapEnum[in.Status]; !ok {
		return in, fmt.Errorf("{\"error\":\"status must be one of [%s, %s, %s]\"}", enum[0], enum[1], enum[2])
	}

	Age, err := strconv.Atoi(r.FormValue("age"))
	if err != nil {
		return in, fmt.Errorf("{\"error\":\"age must be int\"}")
	}
	in.Age = Age
	if in.Age < 0 {
		return in, fmt.Errorf("{\"error\":\"age must be >= 0\"}")
	}
	if in.Age > 128 {
		return in, fmt.Errorf("{\"error\":\"age must be <= 128\"}")
	}
	return in, nil
}

func ValidOtherCreateParams(r *http.Request) (OtherCreateParams, error) {
	in := OtherCreateParams{}

	in.Username = r.FormValue("username")
	if in.Username == "" {
		return in, fmt.Errorf("{\"error\":\"username must me not empty\"}")
	}
	if len(in.Username) < 3 {
		return in, fmt.Errorf("{\"error\":\"username len must be >= 3\"}")
	}

	in.Name = r.FormValue("account_name")

	in.Class = r.FormValue("class")
	if in.Class == "" {
		in.Class = "warrior"
	}
	enum := strings.Split("warrior|sorcerer|rouge", "|")
	mapEnum := make(map[string]string)
	for _, v := range enum {
		mapEnum[v] = ""
	}
	if _, ok := mapEnum[in.Class]; !ok {
		return in, fmt.Errorf("{\"error\":\"class must be one of [%s, %s, %s]\"}", enum[0], enum[1], enum[2])
	}

	Level, err := strconv.Atoi(r.FormValue("level"))
	if err != nil {
		return in, fmt.Errorf("{\"error\":\"level must be int\"}")
	}
	in.Level = Level
	if in.Level < 1 {
		return in, fmt.Errorf("{\"error\":\"level must be >= 1\"}")
	}
	if in.Level > 50 {
		return in, fmt.Errorf("{\"error\":\"level must be <= 50\"}")
	}
	return in, nil
}

func (h *MyApi) ProfileUser(w http.ResponseWriter, r *http.Request) {

	ctx := context.TODO()
	in, err := ValidProfileParams(r)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	res, err := h.Profile(ctx, in)
	if err != nil {
		ae, ok := err.(ApiError)
		if !ok {
			http.Error(w, fmt.Sprintf("{\"error\":\"%s\"}", err.Error()), http.StatusInternalServerError)
			return
		}
		http.Error(w, fmt.Sprintf("{\"error\":\"%s\"}", err.Error()), ae.HTTPStatus)
		return
	}
	result := Result{
		"error": "",
		"response": Result{

			"id":        res.ID,
			"login":     res.Login,
			"full_name": res.FullName,
			"status":    res.Status,
		},
	}
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		http.Error(w, "error encoding JSON", http.StatusInternalServerError)
	}

}
func (h *MyApi) CreateUser(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		password := r.Header.Get("X-Auth")
		if password == "" {
			http.Error(w, "{\"error\":\"unauthorized\"}", 403)
			return
		}
		if password != "100500" {
			http.Error(w, "{\"error\":\"wrong password\"}", http.StatusUnauthorized)
			return
		}
		ctx := context.TODO()
		in, err := ValidCreateParams(r)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		res, err := h.Create(ctx, in)
		if err != nil {
			ae, ok := err.(ApiError)
			if !ok {
				http.Error(w, fmt.Sprintf("{\"error\":\"%s\"}", err.Error()), http.StatusInternalServerError)
				return
			}
			http.Error(w, fmt.Sprintf("{\"error\":\"%s\"}", err.Error()), ae.HTTPStatus)
			return
		}
		result := Result{
			"error": "",
			"response": Result{

				"id": res.ID,
			},
		}
		err = json.NewEncoder(w).Encode(result)
		if err != nil {
			http.Error(w, "error encoding JSON", http.StatusInternalServerError)
		}
	} else {
		http.Error(w, "{\"error\":\"bad method\"}", 406)
	}
}
func (h *OtherApi) CreateUser(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		password := r.Header.Get("X-Auth")
		if password == "" {
			http.Error(w, "{\"error\":\"unauthorized\"}", 403)
			return
		}
		if password != "100500" {
			http.Error(w, "{\"error\":\"wrong password\"}", http.StatusUnauthorized)
			return
		}
		ctx := context.TODO()
		in, err := ValidOtherCreateParams(r)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		res, err := h.Create(ctx, in)
		if err != nil {
			ae, ok := err.(ApiError)
			if !ok {
				http.Error(w, fmt.Sprintf("{\"error\":\"%s\"}", err.Error()), http.StatusInternalServerError)
				return
			}
			http.Error(w, fmt.Sprintf("{\"error\":\"%s\"}", err.Error()), ae.HTTPStatus)
			return
		}
		result := Result{
			"error": "",
			"response": Result{

				"id":        res.ID,
				"login":     res.Login,
				"full_name": res.FullName,
				"level":     res.Level,
			},
		}
		err = json.NewEncoder(w).Encode(result)
		if err != nil {
			http.Error(w, "error encoding JSON", http.StatusInternalServerError)
		}
	} else {
		http.Error(w, "{\"error\":\"bad method\"}", 406)
	}
}
