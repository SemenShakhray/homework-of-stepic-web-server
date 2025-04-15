package handlers

import (
	"fmt"
	"log"
	"net/http"
	token "rwa/internal/lib"
	"rwa/internal/models"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Register(c *gin.Context) {
	var req models.RequestNewUser

	if err := c.ShouldBindJSON(&req); err != nil {
		h.ErrorResponse(c, err, http.StatusBadRequest, "wrong request of registration")

		return
	}

	resp, err := h.service.Register(req)
	if err != nil {
		h.ErrorResponse(c, err, http.StatusInternalServerError, "failed regisration")

		return
	}

	h.responseOK(c, resp, http.StatusCreated)
}

func (h *Handler) Login(c *gin.Context) {
	var req models.RequestLogin

	if err := c.BindJSON(&req); err != nil {
		h.ErrorResponse(c, err, http.StatusBadRequest, "failed request")

		return
	}

	token, err := token.CreateJWT(req.User.Email, h.cfg.TokenTTL)
	if err != nil {
		h.ErrorResponse(c, err, http.StatusInternalServerError, "failed to create token")

		return
	}

	user, err := h.service.Login(req, token)
	if err != nil {
		h.ErrorResponse(c, err, http.StatusInternalServerError, "failed to login")

		return
	}

	log.Println("user after login:", user)
	h.responseOK(c, user, http.StatusOK)

}

func (h *Handler) GetProfile(c *gin.Context) {
	email := c.GetString("email")
	token := c.GetString("token")

	profile, err := h.service.GetUser(email)
	if err != nil {
		h.ErrorResponse(c, err, http.StatusInternalServerError, "failed to get profile")
	}

	if token != profile.User.Token {
		h.ErrorResponse(c, fmt.Errorf("invalid token"), http.StatusBadRequest, "invalid token")
	}

	h.responseOK(c, profile, http.StatusOK)
}

func (h *Handler) UpdateProfile(c *gin.Context) {
	// email := c.GetString("email")

	// profile
}
