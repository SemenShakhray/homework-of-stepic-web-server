package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewRouter(h *http.Handler) *gin.Engine {
	r := gin.Default()

	r.POST("api/users")

	return r
}
