package service

import (
	"rwa/internal/handlers"
	"rwa/internal/models"
)

type Service struct {
	store Storer
}

type Storer interface {
	UsersStore
	ArticleStore
}

type UsersStore interface {
	Register(email, username, password string) error
	AddToken(token, email string) error
	GetUser(email string) (models.User, error)
	UpdateUser(user map[string]string, email string) error
}

type ArticleStore interface {
	CreateWithResponse(req models.RequestNewArticle, slug, email string) (models.Article, error)
}

func NewService(store Storer) handlers.Servicer {
	return &Service{
		store: store,
	}
}
