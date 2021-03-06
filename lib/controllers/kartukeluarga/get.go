package kartukeluarga

import (
	"binadesa2020-backend/lib/clog"
	"binadesa2020-backend/lib/models"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/gin-gonic/gin"
)

// Get kartu-keluarga submission by Admin
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
		var karkel models.KartuKeluarga

		found, _ := karkel.GetByUniqueCode(*req.UniqueCode)
		if found != true {
			c.JSON(http.StatusNotFound, gin.H{"message": "data not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "OK", "data": karkel})
		return
	}

	// Get multiple
	var (
		findOpt     options.FindOptions
		karkelMdl   models.KartuKeluarga
		multikarkel []*models.KartuKeluarga = make([]*models.KartuKeluarga, 0)
	)

	findOpt.SetSort(bson.M{"_id": -1}) // sort by latest ID

	// filtering
	filter := bson.D{}

	// filter by status code
	if req.StatusCode != nil {
		filter = append(filter, bson.E{"status_code", *req.StatusCode})
	}

	cur, err := karkelMdl.Collection().Find(c, filter, &findOpt)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		clog.Fatal(err, "get all kartu keluarga submission")
		return
	}

	for cur.Next(c) {
		var karkel models.KartuKeluarga
		cur.Decode(&karkel)
		multikarkel = append(multikarkel, &karkel)
	}

	c.JSON(http.StatusOK, gin.H{"message": "OK", "data": multikarkel})
}

// GetOne by user
func GetOne(c *gin.Context) {
	var req struct {
		UniqueCode string `form:"unique_code" binding:"required"`
	}

	// extract parameter query
	if err := c.BindQuery(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var karkel models.KartuKeluarga
	found, _ := karkel.GetByUniqueCode(req.UniqueCode)
	if found != true {
		c.JSON(http.StatusNotFound, gin.H{"message": "data not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "OK", "data": karkel})
}
