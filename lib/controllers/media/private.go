package media

import (
	"binadesa2020-backend/lib/services/storage"
	"binadesa2020-backend/lib/variable"
	"bytes"
	"net/http"
	"path"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
)

// DownloadPrivateFile from minio storage for admin
func DownloadPrivateFile(c *gin.Context) {
	var req struct {
		ObjectName string `json:"object_name" binding:"required"`
	}

	// extract parameter
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	// get object from minio storage
	objReader, err := storage.MinioClient.GetObject(
		c,
		variable.ProjectName,
		req.ObjectName,
		minio.GetObjectOptions{},
	)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	// reading buffer from minio reader, and handle if data is empty
	buf := new(bytes.Buffer)
	buf.ReadFrom(objReader)
	size := len(buf.Bytes())
	if size == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "data not found"})
		return
	}

	// get object property
	_, filename := path.Split(req.ObjectName)
	contentType := http.DetectContentType(buf.Bytes())

	// write to response
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", contentType)
	c.Data(http.StatusOK, contentType, buf.Bytes())
}
