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

	resp, err := h.service.Create(req, email)
	if err != nil {
		h.ErrorResponse(c, err, http.StatusInternalServerError, "failed to create article")

		return
	}

	h.responseOK(c, resp, http.StatusCreated)
}
