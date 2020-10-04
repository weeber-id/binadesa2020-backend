package administrator

import (
	"binadesa2020-backend/lib/middleware"
	"binadesa2020-backend/lib/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Delete admin account
func Delete(c *gin.Context) {
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
