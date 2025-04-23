package models

import "time"

type User struct {
	UserID       int       `json:"-"`
	Username     string    `json:"username,omitempty"`
	Email        string    `json:"email,omitempty"`
	PasswordHash string    `json:"-"`
	Bio          string    `json:"bio,omitempty"`
	Image        string    `json:"image,omitempty"`
	Token        string    `json:"token,omitempty"`
	CreatedAt    time.Time `json:"created_at,omitempty"`
	UpdatedAt    time.Time `json:"updated_at,omitempty"`
}

type ResponseUser struct {
	User User `json:"user"`
}

type RequestNewUser struct {
	User struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}
}

type RequestLogin struct {
	User struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}
}

type Author struct {
	Username  string `json:"username"`
	Bio       string `json:"bio"`
	Image     string `json:"image"`
	Following bool   `json:"following"`
}
