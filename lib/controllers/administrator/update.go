package administrator

import (
	"binadesa2020-backend/lib/middleware"
	"binadesa2020-backend/lib/models"
	"binadesa2020-backend/lib/tools"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ChangePassword admin account
func ChangePassword(c *gin.Context) {
	var req struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required"`
	}

	// extract JSON parameter
	if err := c.BindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	// get claims from JWT
	claims := middleware.GetClaims(c)

	admin := &models.Admin{}
	admin.FindByUsername(claims.Username)

	// check old password
	if admin.Password != tools.EncodeMD5(req.OldPassword) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "wrong old password"})
		return
	}

	// change password
	admin.Password = tools.EncodeMD5(req.NewPassword)
	admin.Update()
	c.JSON(http.StatusOK, gin.H{"message": "account updated"})
}
