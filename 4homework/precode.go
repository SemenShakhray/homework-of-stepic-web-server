package main

import (
	"encoding/xml"
	"fmt"
	"os"
	"sort"
	"strings"
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
	Limit      int
	Offset     int    // Можно учесть после сортировки
	Query      string // подстрока в 1 из полей
	OrderField string
	OrderBy    int
}

const (
	OrderByAsc  = -1
	OrderByAsIs = 0
	OrderByDesc = 1

	ErrorBadOrderField = `OrderField invalid`
)

func SearchServer(val SearchRequest) {

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

	searchUsers := make([]User, 0)
	if val.Query == "" {
		searchUsers = append(searchUsers, users.Row...)
	} else {
		for _, user := range users.Row {
			name := user.FirstName + user.LastName
			searchName := strings.Contains(name, val.Query)
			searchAbout := strings.Contains(user.About, val.Query)
			if searchName || searchAbout {
				searchUsers = append(searchUsers, user)
			}
		}
	}
	if len(searchUsers) == 0 {
		fmt.Printf("error: This substring does not exist\n")
		return
	} else {
		switch strings.ToLower(val.OrderField) {
		case "":
			switch val.OrderBy {
			case OrderByDesc:
				sort.Slice(searchUsers, func(i, j int) bool { return searchUsers[i].FirstName < searchUsers[j].FirstName })
			case OrderByAsIs:
			case OrderByAsc:
				sort.Slice(searchUsers, func(i, j int) bool { return searchUsers[j].FirstName < searchUsers[i].FirstName })
			default:
				fmt.Printf("error: %s\n", ErrorBadOrderField)
				return
			}
		case "id":
			switch val.OrderBy {
			case OrderByDesc:
				sort.Slice(searchUsers, func(i, j int) bool { return searchUsers[i].Id < searchUsers[j].Id })
			case OrderByAsIs:
			case OrderByAsc:
				sort.Slice(searchUsers, func(i, j int) bool { return searchUsers[j].Id < searchUsers[i].Id })
			default:
				fmt.Printf("error: %s\n", ErrorBadOrderField)
				return
			}
		case "age":
			switch val.OrderBy {
			case OrderByDesc:
				sort.Slice(searchUsers, func(i, j int) bool { return searchUsers[i].Age < searchUsers[j].Age })
			case OrderByAsIs:
			case OrderByAsc:
				sort.Slice(searchUsers, func(i, j int) bool { return searchUsers[j].Age < searchUsers[i].Age })
			default:
				fmt.Printf("error: %s\n", ErrorBadOrderField)
				return
			}
		default:
			fmt.Printf("error: incorrect sorting parameter\n")
			return
		}
	}
	if val.Limit < 0 {
		fmt.Printf("limit must be > 0\n")
		return
	}
	if val.Limit > 25 {
		val.Limit = 25
	}
	if val.Offset < 0 {
		fmt.Printf("offset must be > 0\n")
		return
	}
	result := searchUsers[val.Offset:]
	for i, v := range result {
		if i > val.Limit-1 {
			continue
		}
		fmt.Printf("Id: %d, Name: %s, Age: %d\n", v.Id, v.FirstName+v.LastName, v.Age)
	}
}

func main() {
	SearchServer(SearchRequest{
		Query:      "",
		OrderField: "AGE",
		OrderBy:    0,
		Limit:      10,
		Offset:     12,
	})
}
