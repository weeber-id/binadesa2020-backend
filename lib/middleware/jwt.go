package middleware

import (
	"binadesa2020-backend/lib/models"
	"binadesa2020-backend/lib/variable"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"

	"github.com/gin-gonic/gin"
)

// AdminAuthorization using JWT
func AdminAuthorization() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get JWT from Header
		authHeader := c.Request.Header["Authorization"]
		if authHeader == nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Missing Authorization Header"})
			c.Abort()
			return
		}
		token := strings.Split(authHeader[0], " ")[1]

		// Decode JWT
		claims := jwt.MapClaims{}
		_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(variable.JWTConfig.Key), nil
		})
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid Access Token"})
			c.Abort()
			return
		}

		// Set Claims to gin context
		c.Set("JWT_ROLE", claims["role"])
		c.Set("JWT_USERNAME", claims["username"])
		c.Set("JWT_NAME", claims["name"])
		c.Set("JWT_LEVEL", int(claims["level"].(float64)))

		c.Next()
	}
}

// GetClaims from gin context and parsed to custom JWTClaims
func GetClaims(c *gin.Context) models.JwtClaims {
	var claims models.JwtClaims
	claims.Role = c.GetString("JWT_ROLE")
	claims.Username = c.GetString("JWT_USERNAME")
	claims.Name = c.GetString("JWT_NAME")
	claims.Level = c.GetInt("JWT_LEVEL")

	return claims
}
