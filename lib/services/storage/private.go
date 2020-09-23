package storage

import (
	"binadesa2020-backend/lib/variable"
	"bytes"
	"context"
	"fmt"

	"github.com/minio/minio-go/v7"
)

// BaseObject minio storage
type BaseObject struct {
	ObjectName string
	File       []byte
}

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
		int64(len(p.File)),
		minio.PutObjectOptions{},
	)
	if err != nil {
		return nil, err
	}

	fmt.Println(info.Location)
	return &info, nil
}

// Download file
func (PrivateObject) Download() error {
	return nil
}

// Delete file
func (PrivateObject) Delete() error {
	return nil
}
