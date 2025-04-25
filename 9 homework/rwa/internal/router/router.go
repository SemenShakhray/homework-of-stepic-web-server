package router

import (
	"rwa/internal/handlers"
	"rwa/internal/handlers/middleware"

	"github.com/gin-gonic/gin"
)

func NewRouter(h *handlers.Handler) *gin.Engine {
	r := gin.Default()

	r.POST("api/users", h.Register)
	r.POST("api/users/login", h.Login)
	r.GET("api/articles", h.GetAllArticlesByFiltres)

	auth := r.Group("api/", middleware.ValidToken(h.Cfg.SecretJWT))
	{
		auth.GET("/user", h.GetUser)
		auth.PUT("/user", h.UpdateUser)
		auth.POST("/articles", h.CreateArticle)
	}

	return r
}
