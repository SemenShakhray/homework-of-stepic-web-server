package handlers

import (
	"fmt"
	"net/http"
	"rwa/internal/lib"
	"rwa/internal/models"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Register(c *gin.Context) {
	var req models.RequestNewUser

	if err := c.ShouldBindJSON(&req); err != nil {
		h.ErrorResponse(c, err, http.StatusBadRequest, "wrong request of registration")

		return
	}

	token, err := lib.CreateJWT(req.User.Email, h.Cfg.TokenTTL)
	if err != nil {
		h.ErrorResponse(c, err, http.StatusInternalServerError, "failed to create token")

		return
	}

	resp, err := h.service.Register(req, token)
	if err != nil {
		h.ErrorResponse(c, err, http.StatusInternalServerError, "failed regisration")

		return
	}

	h.responseOK(c, models.ResponseUser{User: resp}, http.StatusCreated)
}

func (h *Handler) Login(c *gin.Context) {
	var req models.RequestLogin

	if err := c.BindJSON(&req); err != nil {
		h.ErrorResponse(c, err, http.StatusBadRequest, "failed request")

		return
	}

	user, err := h.service.Login(req)
	if err != nil {
		h.ErrorResponse(c, err, http.StatusInternalServerError, "failed to login")

		return
	}

	h.responseOK(c, models.ResponseUser{User: user}, http.StatusOK)
}

func (h *Handler) GetUser(c *gin.Context) {
	email := c.GetString("email")
	token := c.GetString("token")

	user, err := h.service.GetUser(email)
	if err != nil {
		h.ErrorResponse(c, err, http.StatusInternalServerError, "failed to get profile")
	}

	if token != user.Token {
		h.ErrorResponse(c, fmt.Errorf("invalid token"), http.StatusBadRequest, "invalid token")
	}

	h.responseOK(c, models.ResponseUser{User: user}, http.StatusOK)
}

func (h *Handler) UpdateUser(c *gin.Context) {
	var user models.ResponseUser

	err := c.ShouldBindJSON(&user)
	if err != nil {
		h.ErrorResponse(c, err, http.StatusBadRequest, "failed request deserialization ")
	}

	email := c.GetString("email")

	userResp, err := h.service.UpdateUser(user.User, email, h.Cfg.TokenTTL)
	if err != nil {
		h.ErrorResponse(c, err, http.StatusInternalServerError, "failed to update user")

		return
	}

	h.responseOK(c, models.ResponseUser{User: userResp}, http.StatusOK)
}
