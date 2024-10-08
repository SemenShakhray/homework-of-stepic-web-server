package main

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
)

type UserXML struct {
	Id        int    `xml:"id"`
	FirstName string `xml:"first_name"`
	LastName  string `xml:"last_name"`
	Age       int    `xml:"age"`
	About     string `xml:"about"`
	Gender    string `xml:"gender"`
}

type UsersXML struct {
	Row []UserXML `xml:"row"`
}

const key string = "secret"

func EncodingJSON(w http.ResponseWriter, error SearchErrorResponse) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	err := json.NewEncoder(w).Encode(error)
	if err != nil {
		http.Error(w, "error encoding JSON", http.StatusInternalServerError)
	}
}

func SearchServer(w http.ResponseWriter, r *http.Request) {

	auth := r.Header.Get("AccessToken")
	if auth != key {
		http.Error(w, "bad AccessToken", http.StatusUnauthorized)
		return
	}

	read, err := os.ReadFile("dataset.xml")
	if err != nil {
		http.Error(w, "error read file", http.StatusInternalServerError)
		return
	}
	var valueOrderBy, valueLimit, valueOffset int
	users := &UsersXML{}
	err = xml.Unmarshal([]byte(read), users)
	if err != nil {
		http.Error(w, "error parsing XML", http.StatusInternalServerError)
		return
	}
	if r.FormValue("order_by") == "" {
		valueOrderBy = 0
	} else {
		valueOrderBy, err = strconv.Atoi(r.FormValue("order_by"))
		if err != nil {
			http.Error(w, "error parsing to int order_by", http.StatusBadRequest)
		}
	}
	if r.FormValue("limit") == "" {
		valueLimit = 10
	} else {
		valueLimit, err = strconv.Atoi(r.FormValue("limit"))
		if err != nil {
			http.Error(w, "error parsing to int limit", http.StatusBadRequest)
		}
	}
	if r.FormValue("offset") == "" {
		valueOffset = 0
	} else {
		valueOffset, err = strconv.Atoi(r.FormValue("offset"))
		if err != nil {
			http.Error(w, "error parsing to int offset", http.StatusBadRequest)
		}
	}

	searchReq := SearchRequest{
		Query:      r.FormValue("query"),
		OrderField: r.FormValue("order_field"),
		OrderBy:    valueOrderBy,
		Limit:      valueLimit,
		Offset:     valueOffset,
	}

	searchUsers := make([]UserXML, 0)
	if searchReq.Query == "" {
		searchUsers = append(searchUsers, users.Row...)
	} else {
		for _, user := range users.Row {
			name := strings.ToLower(user.FirstName + user.LastName)
			searchName := strings.Contains(name, strings.ToLower(searchReq.Query))
			searchAbout := strings.Contains(strings.ToLower(user.About), strings.ToLower(searchReq.Query))
			if searchName || searchAbout {
				searchUsers = append(searchUsers, user)
			}
		}
	}

	if len(searchUsers) == 0 {
		http.Error(w, "error: This substring does not exist", http.StatusBadRequest)
		return
	}

	switch strings.ToLower(searchReq.OrderField) {
	case "":
		switch searchReq.OrderBy {
		case OrderByDesc:
			sort.Slice(searchUsers, func(i, j int) bool { return searchUsers[i].FirstName < searchUsers[j].FirstName })
		case OrderByAsIs:
		case OrderByAsc:
			sort.Slice(searchUsers, func(i, j int) bool { return searchUsers[j].FirstName < searchUsers[i].FirstName })
		default:
			error := SearchErrorResponse{
				Error: "wrong order_by",
			}
			EncodingJSON(w, error)
			return
		}
	case "name":
		switch searchReq.OrderBy {
		case OrderByDesc:
			sort.Slice(searchUsers, func(i, j int) bool { return searchUsers[i].FirstName < searchUsers[j].FirstName })
		case OrderByAsIs:
		case OrderByAsc:
			sort.Slice(searchUsers, func(i, j int) bool { return searchUsers[j].FirstName < searchUsers[i].FirstName })
		default:
			error := SearchErrorResponse{
				Error: "wrong order_by",
			}
			EncodingJSON(w, error)
			return
		}
	case "id":
		switch searchReq.OrderBy {
		case OrderByDesc:
			sort.Slice(searchUsers, func(i, j int) bool { return searchUsers[i].Id < searchUsers[j].Id })
		case OrderByAsIs:
		case OrderByAsc:
			sort.Slice(searchUsers, func(i, j int) bool { return searchUsers[j].Id < searchUsers[i].Id })
		default:
			error := SearchErrorResponse{
				Error: "wrong order_by",
			}
			EncodingJSON(w, error)
			return
		}
	case "age":
		switch searchReq.OrderBy {
		case OrderByDesc:
			sort.Slice(searchUsers, func(i, j int) bool { return searchUsers[i].Age < searchUsers[j].Age })
		case OrderByAsIs:
		case OrderByAsc:
			sort.Slice(searchUsers, func(i, j int) bool { return searchUsers[j].Age < searchUsers[i].Age })
		default:
			error := SearchErrorResponse{
				Error: "wrong order_by",
			}
			EncodingJSON(w, error)
			return
		}
	default:
		error := SearchErrorResponse{
			Error: "ErrorBadOrderField",
		}
		EncodingJSON(w, error)
		return
	}

	if searchReq.Offset > len(searchUsers) {
		error := SearchErrorResponse{
			Error: "Offset more len(searchslice)",
		}
		EncodingJSON(w, error)
		return
	}
	result := searchUsers[searchReq.Offset:]
	response := []User{}
	for i, v := range result {
		if i > searchReq.Limit-1 {
			continue
		}
		name := v.FirstName + v.LastName
		response = append(response, User{
			Id:     v.Id,
			Name:   name,
			Age:    v.Age,
			Gender: v.Gender,
			About:  v.About,
		})
	}

	w.Header().Add("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(&response)
	if err != nil {
		http.Error(w, "error encoding JSON", http.StatusInternalServerError)
	}
}

// func main() {
// 	r := chi.NewRouter()
// 	server := http.Server{
// 		Addr:         ":8080",
// 		Handler:      r,
// 		ReadTimeout:  10 * time.Second,
// 		WriteTimeout: 10 * time.Second,
// 	}
// 	r.Get("/", SearchServer)
// 	fmt.Println("Starting server at :8080")
// 	server.ListenAndServe()
// }
