package kartukeluarga

import (
	"binadesa2020-backend/lib/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Update kartu keluarga by admin
func Update(c *gin.Context) {
	var (
		reqQue struct {
			UniqueCode string `form:"unique_code" binding:"required"`
		}
		reqBody struct {
			Status string `json:"status" binding:"required"`
		}
	)

	// extract query parameter and check
	if err := c.BindQuery(&reqQue); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	// extract json body and check
	if err := c.BindJSON(&reqBody); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}
	statusInt, _ := strconv.Atoi(reqBody.Status)

	var karkel models.KartuKeluarga
	found, _ := karkel.GetByUniqueCode(reqQue.UniqueCode)
	if found == false {
		c.JSON(http.StatusNotFound, gin.H{"message": "data not found"})
		return
	}

	karkel.StatusCode = statusInt
	karkel.Update()
	c.JSON(http.StatusOK, gin.H{"message": "data has been update"})
}
