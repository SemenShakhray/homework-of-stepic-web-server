package main

type Resp struct {
	Response map[string]interface{} `json:"response,omitempty"`
	Error    string                 `json:"error,omitempty"`
}
