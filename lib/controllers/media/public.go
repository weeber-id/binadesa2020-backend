package media

import (
	"binadesa2020-backend/lib/services/storage"
	"net/http"
	"net/url"
	"path"

	"github.com/gin-gonic/gin/binding"

	"github.com/gin-gonic/gin"
)

// UploadPublicFile to minio storages and give the public URL
func UploadPublicFile(c *gin.Context) {
	var req struct {
		FolderName string `form:"folder_name" binding:"required"`
	}

	if err := c.ShouldBindWith(&req, binding.FormMultipart); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "required file"})
		return
	}

	// new public object
	newObject := &storage.PublicObject{}
	location := path.Join(req.FolderName, url.QueryEscape(fileHeader.Filename))
	newObject.LoadFileHeader(fileHeader, location)
	newObject.Upload(c)

	c.JSON(http.StatusOK, gin.H{
		"message": "OK",
		"data": gin.H{
			"url": newObject.URL,
		},
	})
}
