package news

import (
	"binadesa2020-backend/lib/models"
	"math"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/gin-gonic/gin"
)

// Get news for admin and user
// with filtering
func Get(c *gin.Context) {
	var (
		req struct {
			ID   *string `form:"id"`
			Slug *string `form:"slug"`

			// pagination parameter
			Page           *int    `form:"page"` // start from 1
			ContentPerPage *int    `form:"content_per_page"`
			ExcludeID      *string `form:"exclude_id"`
		}
	)

	// extract query parameter
	if err := c.BindQuery(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "invalid type in query parameter"})
		return
	}

	// --------------------------------- get single data ---------------------------------
	if (req.ID != nil) || (req.Slug != nil) {
		var news models.News

		// Get By ID
		if req.ID != nil {
			news.GetByID(*req.ID)
		}

		// Get By Slug
		if req.Slug != nil {
			news.GetBySlugFromURLQuery(*req.Slug)
		}

		// check if news not found
		if (news == models.News{}) {
			c.JSON(http.StatusNotFound, gin.H{"message": "data not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "OK", "data": news})
		return
	}

	// --------------------------------- Get Multiple with pagination ---------------------------------
	// handle if pagination parameter not found
	if (req.Page == nil) && (req.ContentPerPage == nil) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "pagination parameter not found"})
		return
	}

	var (
		newsMdl   models.News
		multiNews []*models.News = make([]*models.News, 0)
	)

	numSkip := (*req.Page - 1) * (*req.ContentPerPage)
	numLimit := *req.ContentPerPage

	// filtering
	filter := bson.D{}
	if req.ExcludeID != nil {
		id, _ := primitive.ObjectIDFromHex(*req.ExcludeID)
		filter = append(filter, bson.E{"_id", bson.M{"$ne": id}})
	}

	// sort by latest and pagination
	opt := options.Find()
	opt.SetSort(bson.M{"_id": -1})
	opt.SetSkip(int64(numSkip))
	opt.SetLimit(int64(numLimit))

	cur, _ := newsMdl.Collection().Find(c, filter, opt)
	for cur.Next(c) {
		var news models.News

		cur.Decode(&news)
		multiNews = append(multiNews, &news)
	}

	// get max num page
	num, _ := newsMdl.Collection().CountDocuments(c, bson.M{})
	maxPage := math.Ceil(float64(num) / float64(*req.ContentPerPage))

	c.JSON(http.StatusOK, gin.H{
		"message":  "OK",
		"data":     multiNews,
		"max_page": maxPage,
	})
}
