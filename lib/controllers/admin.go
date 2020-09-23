package controllers

import (
	"binadesa2020-backend/lib/clog"
	"binadesa2020-backend/lib/middleware"
	"binadesa2020-backend/lib/models"
	"binadesa2020-backend/lib/tools"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
)

// GetAllAdmin list
func GetAllAdmin(c *gin.Context) {
	var (
		adminMdl models.Admin
		data     []*models.Admin
	)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	cur, err := adminMdl.Collection().Find(ctx, bson.M{})
	clog.Panic2Response(c, err, "Find Admin List")

	for cur.Next(ctx) {
		var row models.Admin
		cur.Decode(&row)
		data = append(data, &row)
	}

	c.JSON(http.StatusOK, gin.H{"message": "OK", "data": data})
}

// CreateAdmin account
func CreateAdmin(c *gin.Context) {
	var (
		adminMdl models.Admin
		req      struct {
			Username string `json:"username" binding:"required"`
			Password string `json:"password" binding:"required"`
			Name     string `json:"name" binding:"required"`
			Level    int    `json:"level"`
		}
	)

	err := c.BindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// Get Jwt Claims
	claims := middleware.GetClaims(c)

	// check user admin if exist and create with level higher than request account
	found := adminMdl.FindByUsername(req.Username)
	if found {
		c.JSON(http.StatusBadRequest, gin.H{"message": "username has been exists, try another username"})
		return
	}
	if req.Level < claims.Level {
		c.JSON(http.StatusForbidden, gin.H{"message": "cannot create account with a higher level than you"})
		return
	}

	// write new admin to database
	newAdmin := &models.Admin{
		Username: req.Username,
		Password: tools.EncodeMD5(req.Password),
		Name:     req.Name,
		Level:    req.Level,
	}

	res, err := newAdmin.Create()
	clog.Panic2Response(c, err, "create item")

	c.JSON(http.StatusOK, gin.H{"message": "OK", "data": res})
}

// DeleteAdmin account
func DeleteAdmin(c *gin.Context) {
	var (
		adminMdl models.Admin
		req      struct {
			Username string `form:"username" binding:"required"`
		}
	)

	// Extract parameter
	err := c.BindQuery(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// Get JWT claims
	claims := middleware.GetClaims(c)

	// Check requirement
	found := adminMdl.FindByUsername(req.Username)
	if !found {
		c.JSON(http.StatusNotFound, gin.H{"message": "username not found"})
		return
	}
	if adminMdl.Level < claims.Level {
		c.JSON(http.StatusForbidden, gin.H{"message": "cannot delete an account with a higher level than you"})
		return
	}

	// Delete admin account
	adminMdl.DeleteByUsername(req.Username)
	c.JSON(http.StatusOK, gin.H{"message": "OK"})
}
