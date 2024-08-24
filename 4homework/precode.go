package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
)

type User struct {
	Id        int    `xml:"id"`
	FirstName string `xml:"first_name"`
	LastName  string `xml:"last_name"`
	Age       int    `xml:"age"`
	About     string `xml:"about"`
	Gender    string `xml:"gender"`
}

type Users struct {
	Row []User `xml:"row"`
}

type SearchRequest struct {
	Limit      string
	Offset     string
	Query      string
	OrderField string
	OrderBy    string
}

const (
	OrderByAsc  = -1
	OrderByAsIs = 0
	OrderByDesc = 1

	ErrorBadOrderField = `OrderField invalid`
)

func SearchServer(w http.ResponseWriter, r *http.Request) {
	read, err := os.ReadFile("dataset.xml")
	if err != nil {
		fmt.Errorf("error: read file -%v", err)
		return
	}
	users := &Users{}
	err = xml.Unmarshal([]byte(read), users)
	if err != nil {
		fmt.Errorf("error unmarshal XML - %v", err)
		return
	}

	searchReq := SearchRequest{
		Query:      r.FormValue("query"),
		OrderField: r.FormValue("order_field"),
		OrderBy:    r.FormValue("order_by"),
		Limit:      r.FormValue("limit"),
		Offset:     r.FormValue("offset"),
	}

	searchUsers := make([]User, 0)
	if searchReq.Query == "" {
		searchUsers = append(searchUsers, users.Row...)
	} else {
		for _, user := range users.Row {
			name := user.FirstName + user.LastName
			searchName := strings.Contains(name, searchReq.Query)
			searchAbout := strings.Contains(user.About, searchReq.Query)
			if searchName || searchAbout {
				searchUsers = append(searchUsers, user)
			}
		}
	}
	if len(searchUsers) == 0 {
		http.Error(w, "error: This substring does not exist", http.StatusBadRequest)
		return
	}
	order_by, err := strconv.Atoi(searchReq.OrderBy)
	if err != nil {
		http.Error(w, "error parsing orderBy", http.StatusBadRequest)
	}
	offset, err := strconv.Atoi(searchReq.Offset)
	if err != nil {
		http.Error(w, "error parsing offset", http.StatusBadRequest)
	}
	limit, err := strconv.Atoi(searchReq.Limit)
	if err != nil {
		http.Error(w, "error parsing limit", http.StatusBadRequest)
	}
	switch strings.ToLower(searchReq.OrderField) {
	case "":
		switch order_by {
		case OrderByDesc:
			sort.Slice(searchUsers, func(i, j int) bool { return searchUsers[i].FirstName < searchUsers[j].FirstName })
		case OrderByAsIs:
		case OrderByAsc:
			sort.Slice(searchUsers, func(i, j int) bool { return searchUsers[j].FirstName < searchUsers[i].FirstName })
		default:
			http.Error(w, fmt.Sprintf("error: %s\n", ErrorBadOrderField), http.StatusBadRequest)
			return
		}
	case "name":
		switch order_by {
		case OrderByDesc:
			sort.Slice(searchUsers, func(i, j int) bool { return searchUsers[i].FirstName < searchUsers[j].FirstName })
		case OrderByAsIs:
		case OrderByAsc:
			sort.Slice(searchUsers, func(i, j int) bool { return searchUsers[j].FirstName < searchUsers[i].FirstName })
		default:
			http.Error(w, fmt.Sprintf("error: %s\n", ErrorBadOrderField), http.StatusBadRequest)
			return
		}
	case "id":
		switch order_by {
		case OrderByDesc:
			sort.Slice(searchUsers, func(i, j int) bool { return searchUsers[i].Id < searchUsers[j].Id })
		case OrderByAsIs:
		case OrderByAsc:
			sort.Slice(searchUsers, func(i, j int) bool { return searchUsers[j].Id < searchUsers[i].Id })
		default:
			http.Error(w, fmt.Sprintf("error: %s\n", ErrorBadOrderField), http.StatusBadRequest)
			return
		}
	case "age":
		switch order_by {
		case OrderByDesc:
			sort.Slice(searchUsers, func(i, j int) bool { return searchUsers[i].Age < searchUsers[j].Age })
		case OrderByAsIs:
		case OrderByAsc:
			sort.Slice(searchUsers, func(i, j int) bool { return searchUsers[j].Age < searchUsers[i].Age })
		default:
			http.Error(w, fmt.Sprintf("error: %s\n", ErrorBadOrderField), http.StatusBadRequest)
			return
		}
	default:
		http.Error(w, "error: incorrect sorting parameter", http.StatusBadRequest)
		return
	}

	if limit < 0 {
		http.Error(w, "limit must be > 0", http.StatusBadRequest)
		return
	}
	if limit > 25 {
		limit = 25
	}
	if offset < 0 || offset > len(searchUsers) {
		http.Error(w, "wrong offset", http.StatusBadRequest)
		return
	}

	result := searchUsers[offset:]
	for i, v := range result {
		if i > limit-1 {
			continue
		}
		resp := &User{
			Id:        v.Id,
			FirstName: v.FirstName,
			LastName:  v.LastName,
			Age:       v.Age,
			Gender:    v.Gender,
			About:     v.About,
		}
		w.Header().Add("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(resp)
		if err != nil {
			http.Error(w, "error encoding JSON", http.StatusInternalServerError)
		}
		// fmt.Printf("Id: %d, Name: %s, Age: %d\n", v.Id, v.FirstName+v.LastName, v.Age)
	}
}

func main() {
	r := chi.NewRouter()
	r.Get("/", SearchServer)
	server := http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	fmt.Println("startin server at :8080")
	server.ListenAndServe()
}
