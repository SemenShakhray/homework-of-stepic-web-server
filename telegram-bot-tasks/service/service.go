package service

import (
	"taskbot/internal/handlers"
	"taskbot/internal/models"
)

type Service struct {
	store Storer
}

type Storer interface {
	GetAllTasks() ([]models.Task, error)
	NewTask(description, assignName string, ownerID int) (map[int]string, error)
	Assing(taskID, assignedID int) error
	Unassing(taskID int) error
	Resolve(taskID int) error
	GetMyTasks(myID int) ([]models.Task, error)
	GetOwnTasks(myID int) ([]models.Task, error)
}

func NewService(store Storer) handlers.Servicer {
	return &Service{
		store: store,
	}
}
