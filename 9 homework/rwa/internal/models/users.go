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

// type Profile struct {
// 	Username  string    `json:"username" binding:"required"`
// 	Email     string    `json:"email" binding:"required,email"`
// 	Password  string    `json:"password" binding:"required"`
// 	Bio       string    `json:"bio"`
// 	CreatedAt time.Time `json:"createdAt"`
// 	UpdatedAt time.Time `json:"updatedAt"`
// 	Image     string    `json:"image"`
// 	Token     string    `json:"token"`
// 	Following bool      `json:"following"`
// }

type Session struct {
	UserID int32
	ID     string
}
