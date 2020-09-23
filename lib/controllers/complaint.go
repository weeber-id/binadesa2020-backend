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

// CreateComplaint to database
// Public access
func CreateComplaint(c *gin.Context) {
	var req struct {
		Name      string `json:"name" binding:"required"`
		RT        string `json:"rt" binding:"required"`
		RW        string `json:"rw" binding:"required"`
		Address   string `json:"address" binding:"required"`
		Complaint string `json:"complaint" binding:"required"`
	}

	// extract parameter
	err := c.BindJSON(&req)
	clog.Panic2Response(c, err, "Binding input JSON")

	// parse to model
	newComplaint := &models.Complaint{
		Name:      req.Name,
		RT:        req.RT,
		RW:        req.RW,
		Address:   req.Address,
		Complaint: req.Complaint,
	}

	// write to database
	res, err := newComplaint.Create()
	clog.Panic(err, "input complaint")

	c.JSON(http.StatusOK, gin.H{"message": "OK", "data": res})
}

// GetAllComplaint list
// Only admin can access this
func GetAllComplaint(c *gin.Context) {
	var (
		compMdl models.Complaint
		data    []*models.Complaint
	)

	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	cur, err := compMdl.Collection().Find(ctx, bson.M{})
	clog.Panic(err, "Find all complaint")

	for cur.Next(ctx) {
		var row models.Complaint
		cur.Decode(&row)
		data = append(data, &row)
	}

	c.JSON(http.StatusOK, gin.H{"message": "OK", "data": data})
}
