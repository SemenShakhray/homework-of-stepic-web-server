package models

import "time"

type Profile struct {
	Username  string        `json:"username"`
	Email     string        `json:"email"`
	Password  string        `json:"password"`
	Bio       string        `json:"bio"`
	CreatedAt time.Duration `json:"createdAt"`
	UpdatedAt time.Duration `json:"updatedAt"`
	Image     string        `json:"image"`
	Token     string        `json:"token"`
	Following bool          `json:"following"`
}

type Session struct {
}

type Article struct {
}
