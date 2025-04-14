package handlers

import (
	"log"
	"rwa/internal/storage"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	storage storage.Storer
}

func NewHandler(store storage.Storer) *Handler {
	return &Handler{
		storage: store,
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
