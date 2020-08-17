package controllers

import (
	"binadesa2020-backend/lib/models"
	"binadesa2020-backend/lib/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
)

// GetAdmin list
func GetAdmin(c *gin.Context) {
	var data []*models.Admin

	row := &models.Admin{}
	results := services.MDB.Collection("admin").Find(bson.M{})
	for results.Next(row) {
		data = append(data, row)
	}

	c.JSON(http.StatusOK, gin.H{"message": "OK", "data": data})
}
