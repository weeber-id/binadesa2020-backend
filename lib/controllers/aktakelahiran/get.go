package aktakelahiran

import (
	"binadesa2020-backend/lib/clog"
	"binadesa2020-backend/lib/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Get akta kelahiran submission by Admin
func Get(c *gin.Context) {
	var req struct {
		UniqueCode *string `form:"unique_code"`
		StatusCode *int    `form:"status_code"`
	}

	// extract paramater from query
	if err := c.BindQuery(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	// Get one by unique code
	if req.UniqueCode != nil {
		var akta models.AktaKelahiran

		found, _ := akta.GetByUniqueCode(*req.UniqueCode)
		if found != true {
			c.JSON(http.StatusNotFound, gin.H{"message": "data not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "OK", "data": akta})
		return
	}

	// Get multiple
	var (
		findOpt   options.FindOptions
		aktaMdl   models.AktaKelahiran
		multiakta []*models.AktaKelahiran = make([]*models.AktaKelahiran, 0)
	)

	findOpt.SetSort(bson.M{"_id": -1}) // sort by latest ID

	// filtering
	filter := bson.D{}

	// filter by status code
	if req.StatusCode != nil {
		filter = append(filter, bson.E{"status_code", *req.StatusCode})
	}

	cur, err := aktaMdl.Collection().Find(c, filter, &findOpt)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		clog.Fatal(err, "get all akta kelahiran submission")
		return
	}

	for cur.Next(c) {
		var akta models.AktaKelahiran
		cur.Decode(&akta)
		multiakta = append(multiakta, &akta)
	}

	c.JSON(http.StatusOK, gin.H{"message": "OK", "data": multiakta})
}

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
