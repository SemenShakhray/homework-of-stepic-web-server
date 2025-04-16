package handlers

import (
	"log"
	"rwa/internal/config.go"
	"rwa/internal/models"
	"time"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Servicer
	Cfg     config.Config
}

type Servicer interface {
	Register(req models.RequestNewUser) (models.User, error)
	Login(req models.RequestLogin, token string) (models.User, error)
	GetUser(email string) (models.User, error)
	UpdateUser(user models.User, email string, exp time.Duration) (models.User, error)
}

func NewHandler(service Servicer, cfg config.Config) *Handler {
	return &Handler{
		service: service,
		Cfg:     cfg,
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
