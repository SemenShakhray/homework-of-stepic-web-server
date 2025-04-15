package middleware

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func CheckToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenHeader := c.GetHeader("Authorization")
		if tokenHeader == "" {
			log.Println("token is empty")

			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization token is required"})
			c.Abort()
			return
		}

		if !strings.Contains(tokenHeader, "Token") {
			log.Println("token don't constrains `Token`")

			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization token is wrong"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(tokenHeader, "Token ")

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %s", t.Method.Alg())
			}

			secretKey := "secret"

			return []byte(secretKey), nil
		})

		if err != nil || !token.Valid {
			log.Println("wrong token:", err, tokenString)

			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization token is wrong"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			log.Println("wrong claims")

			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token payload"})
			c.Abort()
			return
		}

		email, ok := claims["email"]
		if !ok {
			log.Println("payload don't contains email")

			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token payload"})
			c.Abort()
			return
		}

		c.Set("email", email)
		c.Set("token", tokenString)

		c.Next()
	}
}
