package main

type Res map[string]interface{}

type Resp struct {
	Response map[string]interface{} `json:"response,omitempty"`
	Error    error                  `json:"error,omitempty"`
}
