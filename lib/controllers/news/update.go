package news

import (
	"binadesa2020-backend/lib/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Update news by admin
func Update(c *gin.Context) {
	var (
		reqQuery struct {
			ID string `form:"id" binding:"required"`
		}
		reqBody struct {
			Author     string `json:"author" binding:"required"`
			Title      string `json:"title" binding:"required"`
			ImageCover string `json:"image_cover" binding:"required"`
			Content    string `json:"content" binding:"required"`
		}
	)

	// extract query parameter
	if err := c.BindQuery(&reqQuery); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	// exxtract json parameter
	if err := c.BindJSON(&reqBody); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	var news models.News

	// check if content is found
	news.GetByID(reqQuery.ID)
	if (news == models.News{}) {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "data not found"})
		return
	}

	// Update internal variable
	news.Author = reqBody.Author
	news.Title = reqBody.Title
	news.ImageCover = reqBody.ImageCover
	news.Content = reqBody.Content

	news.Update()
	c.JSON(http.StatusOK, gin.H{"message": "data updated"})
}
