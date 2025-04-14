package storage

import "rwa/internal/models"

type Storer interface {
	StorerUsers
}

type StorerUsers interface {
	Register(user models.Profile) (models.Profile, error)
}
