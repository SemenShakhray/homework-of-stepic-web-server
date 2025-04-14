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
	var prof models.Profile

	if err := c.ShouldBindJSON(&prof); err != nil {
		h.ErrorResponse(c, err, http.StatusBadRequest, "wrong request of registration")

		return
	}

	resp, err := h.storage.Register(prof)
	if err != nil {
		h.ErrorResponse(c, err, http.StatusInternalServerError, "failed regisration")

		return
	}

	h.responseOK(c, resp, http.StatusCreated)
}

func (h *Handler) Login(c *gin.Context) {
	var req models.Login

	if err := c.BindJSON(&req); err != nil {
		h.ErrorResponse(c, err, http.StatusBadRequest, "failed request")

		return
	}

	pass, err := h.storage.GetPassword(req)
	if err != nil {
		h.ErrorResponse(c, err, http.StatusInternalServerError, "wrong password")

		return
	}

	if req.User.Password != pass {
		h.ErrorResponse(c, fmt.Errorf("wrong password"), http.StatusInternalServerError, "wrong password")

		return
	}

	resp, err := h.storage.Login(req)
	if err != nil {
		h.ErrorResponse(c, err, http.StatusInternalServerError, "failed to login")

		return
	}

	token, err := CreateJWT(resp)
	if err != nil {
		h.ErrorResponse(c, err, http.StatusInternalServerError, "failed to create token")

		return
	}

	resp.User.Token = token

	h.responseOK(c, resp, http.StatusOK)

}

func CreateJWT(profile models.Profile) (string, error) {

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": profile.User.Username,
		"email":    profile.User.Email,
	})

	token, err := claims.SignedString([]byte("secret"))
	if err != nil {
		log.Println("failed to sing token", err)

		return "", err
	}

	return token, nil
}
