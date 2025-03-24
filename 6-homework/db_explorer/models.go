package main

type Res map[string]interface{}

type Resp struct {
	Response interface{} `json:"response,omitempty"`
	Error    string      `json:"error,omitempty"`
}
