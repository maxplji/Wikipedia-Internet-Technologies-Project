package middlewares

import (
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuth() gin.HandlerFunc {
	secret := []byte(os.Getenv("SECRET"))

	return func(c *gin.Context) {
		cookie, err := c.Cookie("Authorization")
		if err != nil || len(secret) == 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Nie jestes zalogowany"})
			return
		}

		claims := &jwt.RegisteredClaims{}
		_, err = jwt.ParseWithClaims(cookie, claims, func(_ *jwt.Token) (any, error) {
			return secret, nil
		}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}), jwt.WithExpirationRequired())
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Nieprawidlowy token"})
			return
		}

		userID, err := strconv.ParseUint(claims.Subject, 10, 64)
		if err != nil || userID == 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Nieprawidlowy token"})
			return
		}

		c.Set("userID", userID)
		c.Next()
	}
}
