package suratketerangan

import (
	"binadesa2020-backend/lib/clog"
	"binadesa2020-backend/lib/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Get all or one for admin
func Get(c *gin.Context) {
	var req struct {
		UniqueCode *string `form:"unique_code"`
		StatusCode *int    `form:"status_code"`
	}

	// extract parameter
	if err := c.BindQuery(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	// get one surat keterangan
	if req.UniqueCode != nil {
		var suratket models.SuratKeterangan
		found, _ := suratket.GetByUniqueCode(*req.UniqueCode)
		if found == false {
			c.JSON(http.StatusNotFound, gin.H{"message": "data not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "OK", "data": suratket})
		return
	}

	// get multiple surat keterangan
	var (
		findOpt       options.FindOptions
		skmdl         models.SuratKeterangan
		multisuratket []*models.SuratKeterangan
	)

	findOpt.SetSort(bson.M{"_id": -1}) // sort by latest

	// filtering
	filter := bson.D{}

	// filter by status code
	if req.StatusCode != nil {
		filter = append(filter, bson.E{"status_code", *req.StatusCode})
	}

	cur, err := skmdl.Collection().Find(c, filter, &findOpt)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		clog.Fatal(err, "get all kartu keluarga submission")
		return
	}

	for cur.Next(c) {
		var suratket models.SuratKeterangan
		cur.Decode(&suratket)
		multisuratket = append(multisuratket, &suratket)
	}

	c.JSON(http.StatusOK, gin.H{"message": "OK", "data": multisuratket})
}

// GetOne for user check progress
func GetOne(c *gin.Context) {
	var req struct {
		UniqueCode string `form:"unique_code" binding:"required"`
	}

	// extract parameter
	if err := c.BindQuery(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	// get one surat keterangan
	var suratket models.SuratKeterangan
	found, _ := suratket.GetByUniqueCode(req.UniqueCode)
	if found == false {
		c.JSON(http.StatusNotFound, gin.H{"message": "data not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "OK", "data": suratket})
	return
}
