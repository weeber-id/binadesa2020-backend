package storage

import (
	"binadesa2020-backend/lib/variable"
	"bytes"
	"context"
	"net/http"
	"path"

	"github.com/minio/minio-go/v7"
)

// PublicObject minio storage
type PublicObject struct {
	BaseObject
	Location string
	URL      string
}

// Upload file to public bucket
// rewrite public URL to this variable
func (p *PublicObject) Upload(ctx context.Context) (*minio.UploadInfo, error) {
	bucket := "public"
	p.ObjectName = path.Join(variable.ProjectName, p.ObjectName) // p.objectname is from controller

	info, err := MinioClient.PutObject(
		ctx,
		"public",
		p.ObjectName,
		bytes.NewReader(p.File),
		p.Size,
		minio.PutObjectOptions{ContentType: http.DetectContentType(p.File)},
	)
	if err != nil {
		return nil, err
	}

	p.URL = path.Join(
		variable.MinioConfig.URIEndpoint,
		bucket,
		p.ObjectName,
	)
	return &info, nil
}
