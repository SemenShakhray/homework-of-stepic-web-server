package token

import (
	"log"
	"time"

	"github.com/golang-jwt/jwt"
)

func CreateJWT(email string, exp time.Duration) (string, error) {

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(exp).Unix(),
		"iat":   time.Now().Unix(),
	})

	token, err := claims.SignedString([]byte("secret"))
	if err != nil {
		log.Println("failed to sing token", err)

		return "", err
	}

	return token, nil
}
