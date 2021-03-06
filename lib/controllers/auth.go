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

	svcConfig := variable.ServiceConfig
	// c.SetCookie(svcConfig.TokenName, token, 3600*24, svcConfig.Path, svcConfig.Domain, svcConfig.HTTPS, true)
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     svcConfig.TokenName,
		Value:    token,
		Path:     svcConfig.Path,
		Domain:   svcConfig.Domain,
		MaxAge:   3600 * 24,
		Secure:   svcConfig.HTTPS,
		HttpOnly: true,
		SameSite: 4, // None
	})

	data := struct {
		Admin       models.Admin `json:"admin"`
		AccessToken string       `json:"access_token"`
	}{
		Admin:       admin,
		AccessToken: token,
	}

	c.JSON(http.StatusOK, gin.H{"data": data})
}

// Logout admin
// delete cookies from backend
func Logout(c *gin.Context) {
	config := variable.ServiceConfig
	// c.SetCookie(config.TokenName, "", 0, config.Path, config.Domain, config.HTTPS, true)
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     config.TokenName,
		Value:    "",
		Path:     config.Path,
		Domain:   config.Domain,
		MaxAge:   0,
		Secure:   config.HTTPS,
		HttpOnly: true,
		SameSite: 4, // None
	})

	c.JSON(http.StatusOK, gin.H{"message": "logout success"})
}
