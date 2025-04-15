package models

import "time"

type Users struct {
	User struct {
		Username     string    `json:"username"`
		Email        string    `json:"email"`
		PasswordHash string    `json:"-"`
		Bio          string    `json:"bio"`
		Image        string    `json:"image"`
		Token        string    `json:"token"`
		CreatedAt    time.Time `json:"createdAt"`
		UpdatedAt    time.Time `json:"updatedAt"`
	}
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

type ResponseUser struct {
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
