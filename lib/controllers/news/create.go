package news

import (
	"binadesa2020-backend/lib/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Create one news by admin
func Create(c *gin.Context) {
	var req struct {
		Author     string `json:"author" binding:"required"`
		Title      string `json:"title" binding:"title"`
		ImageCover string `json:"image_cover" binding:"image_cover"`
		Content    string `json:"content" binding:"content"`
	}

	// extract JSON parameter
	if err := c.BindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	newItem := &models.News{
		Author:     req.Author,
		Title:      req.Title,
		ImageCover: req.ImageCover,
		Content:    req.Content,
	}

	result, err := newItem.Create()
	if err != nil {
		log.Printf("error in create news: %v \n", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "created",
		"data": gin.H{
			"id":   result.InsertedID,
			"slug": newItem.Slug,
		},
	})
}
