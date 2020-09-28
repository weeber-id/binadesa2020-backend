package aktakelahiran

import (
	"binadesa2020-backend/lib/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetOne aktakelahiran submission by user
func GetOne(c *gin.Context) {
	var req struct {
		UniqueCode string `form:"unique_code" binding:"required"`
	}

	// extract parameter query
	if err := c.BindQuery(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var akta models.AktaKelahiran
	found, _ := akta.GetByUniqueCode(req.UniqueCode)
	if found != true {
		c.JSON(http.StatusNotFound, gin.H{"message": "data not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "OK", "data": akta})
}
