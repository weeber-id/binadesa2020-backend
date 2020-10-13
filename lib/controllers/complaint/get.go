package complaint

import (
	"binadesa2020-backend/lib/clog"
	"binadesa2020-backend/lib/models"
	"context"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
)

// Get list
// Only admin can access this
func Get(c *gin.Context) {
	var (
		req struct {
			ID *string `form:"id"`
		}
	)

	// extract parameter
	if err := c.BindQuery(&req); err != nil {
		c.JSON(http.StatusOK, err.Error())
		return
	}

	// --------------------- Get one by ID ---------------------
	if req.ID != nil {
		var comp models.Complaint

		found, _ := comp.GetByID(*req.ID)
		if found == false {
			c.JSON(http.StatusNotFound, gin.H{"message": "data not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "OK", "data": comp})
		return
	}

	// --------------------- Get all by ID ---------------------
	var (
		findOpt options.FindOptions
		compMdl models.Complaint
		data    []*models.Complaint
	)

	findOpt.SetSort(bson.M{"_id": -1})

	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	cur, err := compMdl.Collection().Find(ctx, bson.M{}, &findOpt)
	clog.Panic(err, "Find all complaint")

	for cur.Next(ctx) {
		var row models.Complaint
		cur.Decode(&row)
		data = append(data, &row)
	}

	c.JSON(http.StatusOK, gin.H{"message": "OK", "data": data})
}
