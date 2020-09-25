package kartukeluarga

import (
	"binadesa2020-backend/lib/clog"
	"binadesa2020-backend/lib/models"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"

	"github.com/gin-gonic/gin"
)

// Get kartu-keluarga submission by Admin
func Get(c *gin.Context) {
	var req struct {
		UniqueCode *string `form:"unique_code"`
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
		multikarkel []*models.KartuKeluarga
	)

	findOpt.SetSort(bson.M{"_id": -1}) // sort by latest ID

	cur, err := karkelMdl.Collection().Find(c, bson.M{}, &findOpt)
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
