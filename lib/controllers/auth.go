package controllers

import (
	"binadesa2020-backend/lib/clog"
	"binadesa2020-backend/lib/models"
	"binadesa2020-backend/lib/variable"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// Login admin controller
func Login(c *gin.Context) {
	var (
		admin models.Admin
		req   struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
	)

	err := c.BindJSON(&req)
	clog.Panic2Response(c, err, "Binding JSON")

	ok := admin.Verify(req.Username, req.Password)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid credential"})
		return
	}

	claims := jwt.MapClaims{}
	claims["exp"] = time.Now().Add(24 * time.Hour).Unix()
	claims["role"] = "admin"
	claims["name"] = admin.Name
	claims["username"] = admin.Username
	claims["level"] = admin.Level

	config := variable.JWTConfig
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := at.SignedString([]byte(config.Key))
	clog.Panic2Response(c, err, "generate JWT token")

	data := struct {
		Admin       models.Admin `json:"admin"`
		AccessToken string       `json:"access_token"`
	}{
		Admin:       admin,
		AccessToken: token,
	}

	c.JSON(http.StatusOK, gin.H{"data": data})
}
