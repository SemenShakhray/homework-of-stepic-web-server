package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func testTimeout(w http.ResponseWriter, r *http.Request) {
	time.Sleep(2 * time.Second)
}

func testInternalServerError(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "error of server", http.StatusInternalServerError)
	return
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

	casesRequest := []SearchRequest{
		SearchRequest{
			Query:      "est ea",
			OrderField: "NAME",
			Limit:      26,
			OrderBy:    OrderByAsIs,
			Offset:     0,
		},
		SearchRequest{
			Query:      "est ea",
			OrderField: "NAME",
			Limit:      1,
			OrderBy:    OrderByAsIs,
			Offset:     0,
		},
	}

	for caseNum, item := range casesRequest {
		resp, err := c.FindUsers(item)
		assert.NoError(t, err)
		switch caseNum {
		case 0:
			assert.Equal(t, &SearchResponse{
				Users:    users[:2],
				NextPage: false,
			}, resp)
		case 1:
			assert.Equal(t, &SearchResponse{
				Users:    users[:1],
				NextPage: true,
			}, resp)
		}
	}

	casesErrorClient := []SearchRequest{
		SearchRequest{
			Limit: -1,
		},
		SearchRequest{
			Offset: -1,
		},
		SearchRequest{
			OrderBy: -2,
		},
		SearchRequest{
			OrderField: "Named",
		},
		SearchRequest{
			Query: "Hggasjhgdjasd",
		},
	}

	for _, item := range casesErrorClient {
		_, err := c.FindUsers(item)
		assert.Error(t, err)
	}

	tsTimeout := httptest.NewServer(http.HandlerFunc(testTimeout))
	defer tsTimeout.Close()

	tsInternalServerError := httptest.NewServer(http.HandlerFunc(testInternalServerError))
	defer tsInternalServerError.Close()

	casesClient := []SearchClient{
		SearchClient{
			AccessToken: "sec",
			URL:         ts.URL,
		},

		SearchClient{
			AccessToken: "secret",
			URL:         "",
		},

		SearchClient{
			AccessToken: "",
			URL:         "http://vk.com",
		},

		SearchClient{
			AccessToken: "secret",
			URL:         tsTimeout.URL,
		},

		SearchClient{
			AccessToken: "secret",
			URL:         tsInternalServerError.URL,
		},
	}

	for _, item := range casesClient {
		_, err := item.FindUsers(SearchRequest{})
		assert.Error(t, err)
	}
}
