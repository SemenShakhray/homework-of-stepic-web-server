package handlers

import (
	"fmt"
	"log"
	"net/http"
	"rwa/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
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

	resp, err := h.service.Login(req)
	if err != nil {
		h.ErrorResponse(c, err, http.StatusInternalServerError, "failed to login")

		return
	}

	token, err := CreateJWT(req)
	if err != nil {
		h.ErrorResponse(c, err, http.StatusInternalServerError, "failed to create token")

		return
	}

	err = h.service.AddToken(token, resp.User.Email)
	if err != nil {
		h.ErrorResponse(c, err, http.StatusInternalServerError, "failed to create token")

		return
	}

	resp.User.Token = token

	h.responseOK(c, resp, http.StatusOK)

}

func (h *Handler) GetProfile(c *gin.Context) {
	email := c.GetString("email")
	token := c.GetString("token")

	profile, err := h.storage.GetProfile(email)
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

func CreateJWT(email string) (string, error) {

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
	})

	token, err := claims.SignedString([]byte("secret"))
	if err != nil {
		log.Println("failed to sing token", err)

		return "", err
	}

	return token, nil
}
