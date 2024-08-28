package main

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
	"net/http/httptest"
	"os"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
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

// func SearchRequest
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
	users := &UsersXML{}
	err = xml.Unmarshal([]byte(read), users)
	if err != nil {
		http.Error(w, "error parsing XML", http.StatusInternalServerError)
		return
	}
	valueOrderBy, err := strconv.Atoi(r.FormValue("order_by"))
	if err != nil {
		http.Error(w, "error parsing to int order_by", http.StatusInternalServerError)
	}
	valueLimit, err := strconv.Atoi(r.FormValue("limit"))
	if err != nil {
		http.Error(w, "error parsing to int limit", http.StatusInternalServerError)
	}
	valueOffset, err := strconv.Atoi(r.FormValue("offset"))
	if err != nil {
		http.Error(w, "error parsing to int offset", http.StatusInternalServerError)
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
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			err = json.NewEncoder(w).Encode(error)
			if err != nil {
				http.Error(w, "error encoding JSON", http.StatusInternalServerError)
			}
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
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			err = json.NewEncoder(w).Encode(error)
			if err != nil {
				http.Error(w, "error encoding JSON", http.StatusInternalServerError)
			}
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
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			err = json.NewEncoder(w).Encode(error)
			if err != nil {
				http.Error(w, "error encoding JSON", http.StatusInternalServerError)
			}
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
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			err = json.NewEncoder(w).Encode(error)
			if err != nil {
				http.Error(w, "error encoding JSON", http.StatusInternalServerError)
			}
			return
		}
	default:
		error := SearchErrorResponse{
			Error: "ErrorBadOrderField",
		}
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		err = json.NewEncoder(w).Encode(error)
		if err != nil {
			http.Error(w, "error encoding JSON", http.StatusInternalServerError)
		}
		return
	}

	result := searchUsers[searchReq.Offset:]
	response := []User{}
	for i, v := range result {
		if i > searchReq.Limit-1 {
			return
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
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "error encoding JSON", http.StatusInternalServerError)
	}
}

func TestFindUser(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	defer ts.Close()

	c := SearchClient{
		AccessToken: "secret",
		URL:         ts.URL,
	}

	users := []User{
		User{
			Id:     21,
			Name:   "JohnsWhitney",
			Age:    26,
			About:  "Elit sunt exercitation incididunt est ea quis do ad magna. Commodo laboris nisi aliqua eu incididunt eu irure. Labore ullamco quis deserunt non cupidatat sint aute in incididunt deserunt elit velit. Duis est mollit veniam aliquip. Nulla sunt veniam anim et sint dolore.\n",
			Gender: "male",
		},
		User{
			Id:     28,
			Name:   "CohenHines",
			Age:    32,
			About:  "Deserunt deserunt dolor ex pariatur dolore sunt labore minim deserunt. Tempor non et officia sint culpa quis consectetur pariatur elit sunt. Anim consequat velit exercitation eiusmod aute elit minim velit. Excepteur nulla excepteur duis eiusmod anim reprehenderit officia est ea aliqua nisi deserunt officia eiusmod. Officia enim adipisicing mollit et enim quis magna ea. Officia velit deserunt minim qui. Commodo culpa pariatur eu aliquip voluptate culpa ullamco sit minim laboris fugiat sit.\n",
			Gender: "male",
		},
	}

	checkSearchUsers := SearchRequest{
		Query:      "est ea",
		OrderField: "NAME",
		Limit:      26,
		OrderBy:    OrderByAsIs,
		Offset:     0,
	}
	resp, err := c.FindUsers(checkSearchUsers)
	assert.NoError(t, err)
	assert.Equal(t, &SearchResponse{
		Users:    users[:2],
		NextPage: false,
	}, resp)

	checkSearchNextPage := SearchRequest{
		Query:      "est ea",
		OrderField: "NAME",
		Limit:      1,
		OrderBy:    OrderByAsIs,
		Offset:     0,
	}
	resp, err = c.FindUsers(checkSearchNextPage)
	assert.NoError(t, err)
	assert.Equal(t, &SearchResponse{
		Users:    users[:1],
		NextPage: true,
	}, resp)

	checkLimit := SearchRequest{
		Limit: -1,
	}
	_, err = c.FindUsers(checkLimit)
	assert.Error(t, err)

	checkOffset := SearchRequest{
		Offset: -1,
	}
	_, err = c.FindUsers(checkOffset)
	assert.Error(t, err)

	checkOrderBy := SearchRequest{
		OrderBy: -2,
	}
	_, err = c.FindUsers(checkOrderBy)
	assert.Error(t, err)

	checkOrderField := SearchRequest{
		OrderField: "NAMES",
	}
	_, err = c.FindUsers(checkOrderField)
	assert.Error(t, err)

	checkCantUnpackJSON := SearchRequest{
		Query: "dsfdsfsdfsdfsd",
	}
	_, err = c.FindUsers(checkCantUnpackJSON)
	assert.Error(t, err)

	checkAuth := SearchClient{
		AccessToken: "sec",
		URL:         ts.URL,
	}
	_, err = checkAuth.FindUsers(SearchRequest{})
	assert.Error(t, err)

	checkURL := SearchClient{
		AccessToken: "secret",
		URL:         "http://localhost:8081",
	}
	_, err = checkURL.FindUsers(SearchRequest{})
	assert.NoError(t, err)

}
