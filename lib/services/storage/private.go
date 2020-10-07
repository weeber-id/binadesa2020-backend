package storage

import (
	"binadesa2020-backend/lib/variable"
	"bytes"
	"context"
	"fmt"
	"net/http"

	"github.com/minio/minio-go/v7"
)

// PrivateObject minio storage
type PrivateObject struct {
	BaseObject
	Location string
}

// Upload file
func (p *PrivateObject) Upload(ctx context.Context) (*minio.UploadInfo, error) {
	info, err := MinioClient.PutObject(
		ctx,
		variable.ProjectName,
		p.ObjectName,
		bytes.NewReader(p.File),
		p.Size,
		minio.PutObjectOptions{},
	)
	if err != nil {
		return nil, err
	}
	return &info, nil
}

// Download file from objectname
// returns bytes, content-type, error
func (p *PrivateObject) Download(ctx context.Context) ([]byte, string, error) {
	reader, err := MinioClient.GetObject(
		ctx,
		variable.ProjectName,
		p.ObjectName,
		minio.GetObjectOptions{},
	)
	if err != nil {
		return nil, "", err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(reader)
	size := len(buf.Bytes())
	contentType := http.DetectContentType(buf.Bytes())
	if size == 0 {
		return nil, "", fmt.Errorf("file %s not found", p.ObjectName)
	}
	return buf.Bytes(), contentType, nil
}

// Delete file
func (PrivateObject) Delete() error {
	return nil
}
