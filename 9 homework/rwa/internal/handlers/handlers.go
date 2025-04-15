package handlers

import (
	"log"
	"rwa/internal/config.go"
	"rwa/internal/models"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Servicer
	cfg     config.Config
}

type Servicer interface {
	Register(req models.RequestNewUser) (models.Users, error)
	Login(req models.RequestLogin, token string) (models.Users, error)
	GetUser(email string) (models.Users, error)
}

func NewHandler(service Servicer, cfg config.Config) *Handler {
	return &Handler{
		service: service,
		cfg:     cfg,
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
