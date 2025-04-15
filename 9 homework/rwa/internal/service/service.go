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
	// GetPassword(profile models.RequestLogin) (string, error)
	// Login(profile models.RequestLogin) (models.Users, error)
	AddToken(token, email string) error
	GetUser(email string) (models.Users, error)
}

func NewService(store Storer) handlers.Servicer {
	return &Service{
		store: store,
	}
}
