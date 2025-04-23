package handlers

import (
	"log"
	"net/http"
	"rwa/internal/models"

	"github.com/gin-gonic/gin"
)

func (h *Handler) CreateArticle(c *gin.Context) {
	var req models.RequestNewArticle

	email := c.GetString("email")

	err := c.ShouldBindJSON(&req)
	if err != nil {
		log.Println("wrong request of adding new article:", err)

		h.ErrorResponse(c, err, http.StatusBadRequest, "wrong request of registration")
		return
	}

	res, err := h.service.Create(req, email)
	if err != nil {
		h.ErrorResponse(c, err, http.StatusInternalServerError, "failed to create article")

		return
	}

	resp := models.Articles{
		Article: res,
	}

	h.responseOK(c, resp, http.StatusCreated)
}

func (h *Handler) GetAllArticlesByFiltres(c *gin.Context) {
	var params models.ArticleQueryParams

	if err := c.ShouldBindQuery(&params); err != nil {
		h.ErrorResponse(c, err, http.StatusBadRequest, "invalid query parameters")
	}
}
