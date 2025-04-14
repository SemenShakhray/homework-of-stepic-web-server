package models

import "time"

type Profile struct {
	User struct {
		Username  string    `json:"username" binding:"required"`
		Email     string    `json:"email" binding:"required,email"`
		Password  string    `json:"password" binding:"required"`
		Bio       string    `json:"bio"`
		CreatedAt time.Time `json:"createdAt"`
		UpdatedAt time.Time `json:"updatedAt"`
		Image     string    `json:"image"`
		Token     string    `json:"token"`
		Following bool      `json:"following"`
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
}

type Article struct {
}
