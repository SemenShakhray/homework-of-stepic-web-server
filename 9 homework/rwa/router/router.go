package router

import (
	"rwa/internal/handlers"

	"github.com/gin-gonic/gin"
)

func NewRouter(h *handlers.Handler) *gin.Engine {
	r := gin.Default()

	r.POST("api/users", h.Register)
	r.POST("api/users/login", h.Login)

	return r
}
