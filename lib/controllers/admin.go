package controllers

import (
	"binadesa2020-backend/lib/clog"
	"binadesa2020-backend/lib/models"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
)

// GetAdmin list
func GetAdmin(c *gin.Context) {
	var (
		adminMdl models.Admin
		data     []*models.Admin
	)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	collection := adminMdl.Collection()

	cur, err := collection.Find(ctx, bson.M{})
	clog.Panic2Response(c, err, "Find Admin List")

	for cur.Next(ctx) {
		var row models.Admin
		cur.Decode(&row)
		data = append(data, &row)
	}

	c.JSON(http.StatusOK, gin.H{"message": "OK", "data": data})
}
