package handlers

import (
	"log"
	"rwa/internal/models"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Servicer
}

type Servicer interface {
	Register(req models.RequestNewUser) (models.Users, error)
	Login(req models.RequestLogin) (models.Users, error)
	GetUser(email string) (models.Users, error)
}

func NewHandler(service Servicer) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) ErrorResponse(c *gin.Context, err error, code int, message string) {
	log.Println(message, ":", err)

	c.JSON(code, gin.H{
		"error": message,
	})
}

func (h *Handler) responseOK(c *gin.Context, resp any, code int) {
	c.JSON(code, resp)
}
