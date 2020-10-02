package news

import (
	"binadesa2020-backend/lib/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Delete news by admin
func Delete(c *gin.Context) {
	var req struct {
		ID string `form:"id" binding:"required"`
	}

	// extract parameter
	if err := c.BindQuery(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	var news models.News

	news.GetByID(req.ID)
	if (news == models.News{}) {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "data not found"})
		return
	}

	news.Delete()
	c.JSON(http.StatusOK, gin.H{"message": "data deleted"})
}
