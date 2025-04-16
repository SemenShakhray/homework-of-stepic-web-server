package service

import (
	"rwa/internal/handlers"
	"rwa/internal/models"
)

type Service struct {
	store Storer
}

type Storer interface {
	StorerUsers
}

type StorerUsers interface {
	Register(email, username, password string) error
	AddToken(token, email string) error
	GetUser(email string) (models.User, error)
	UpdateUser(user map[string]string, email string) error
}

func NewService(store Storer) handlers.Servicer {
	return &Service{
		store: store,
	}
}
