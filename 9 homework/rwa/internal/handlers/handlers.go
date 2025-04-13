package handlers

import "rwa/internal/storage"

type Handler struct {
	storage storage.Storer
}

func NewHandler(store storage.Storer) *Handler {
	return &Handler{
		storage: store,
	}
}
