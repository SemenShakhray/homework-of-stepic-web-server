package storage

import "rwa/internal/models"

type Storer interface {
	StorerUsers
}

type StorerUsers interface {
	Register(profile models.Profile) (models.Profile, error)
}
