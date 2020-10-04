package complaint

import (
	"binadesa2020-backend/lib/clog"
	"binadesa2020-backend/lib/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Create to database
// Public access
func Create(c *gin.Context) {
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
