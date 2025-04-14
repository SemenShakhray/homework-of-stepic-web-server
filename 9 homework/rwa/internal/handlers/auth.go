package handlers

import (
	"log"
	"net/http"
	"rwa/internal/models"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Register(c *gin.Context) {
	var prof models.Profile

	if err := c.ShouldBindJSON(&prof); err != nil {
		h.ErrorResponse(c, err, http.StatusBadRequest, "wrong request of registration")
	}

	resp, err := h.storage.Register(prof)
	if err != nil {
		h.ErrorResponse(c, err, http.StatusInternalServerError, "failed regisration")
	}

	log.Println(resp)
	h.responseOK(c, resp, http.StatusCreated)
}
