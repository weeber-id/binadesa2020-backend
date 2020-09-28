package controllers

import (
	"binadesa2020-backend/lib/variable"
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthCheck with message and datetime now
func HealthCheck(c *gin.Context) {
	now := variable.DateTimeNowPtr()
	c.JSON(http.StatusOK, gin.H{
		"message":  "works",
		"version":  variable.Version,
		"datetime": now,
	})
}
