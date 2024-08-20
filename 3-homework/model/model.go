package model

type User struct {
	Browsers []string `json:"browsers"`
	Email    string   `json:"email,omitempty"`
	Name     string   `json:"name,omitempty"`
}
