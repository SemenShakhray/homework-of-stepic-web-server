package storage

import "rwa/internal/models"

type Storer interface {
	StorerUsers
}

type StorerUsers interface {
	Register(user models.Profile) (models.Profile, error)
	GetPassword(profile models.Login) (string, error)
	Login(profile models.Login) (models.Profile, error)
}
