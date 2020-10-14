package suratketerangan

import (
	"binadesa2020-backend/lib/models"
	"binadesa2020-backend/lib/services/gmail"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Update surat keterangan submission by admin
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

	var suratket models.SuratKeterangan
	found, _ := suratket.GetByUniqueCode(reqQue.UniqueCode)
	if found == false {
		c.JSON(http.StatusNotFound, gin.H{"message": "data not found"})
		return
	}

	suratket.StatusCode = statusInt
	suratket.Update()

	// Send Email Status
	if statusInt == models.StatusCode.Accepted {
		go func() {
			email := gmail.Email{To: suratket.Email}
			email.SendCompleteSubmission(&suratket)
		}()
	} else if statusInt == models.StatusCode.Rejected {
		go func() {
			email := gmail.Email{To: suratket.Email}
			email.SendRejectSubmission(&suratket)
		}()
	}

	c.JSON(http.StatusOK, gin.H{"message": "data has been update"})
}
