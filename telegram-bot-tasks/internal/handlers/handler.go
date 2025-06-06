package handlers

type Servicer interface {
	GetAllTasks(userID int, userName string) (map[int]string, error)
	NewTask(description, assingName string, ownerID int) (map[int]string, error)
}

type Handler struct {
	serv Servicer
}

func NewHandler(serv Servicer) *Handler {
	return &Handler{
		serv: serv,
	}
}
