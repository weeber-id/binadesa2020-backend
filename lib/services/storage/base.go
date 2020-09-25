package storage

import (
	"mime/multipart"
)

// BaseObject minio storage
type BaseObject struct {
	ObjectName string
	File       []byte
	Size       int64
}

// LoadFileHeader to this struct
// Usually used from gin formdata
func (b *BaseObject) LoadFileHeader(fileHeader *multipart.FileHeader, objectName string) error {
	b.ObjectName = objectName
	b.Size = fileHeader.Size

	file, err := fileHeader.Open()
	defer file.Close()
	if err != nil {
		return err
	}

	buf := make([]byte, fileHeader.Size)
	_, err = file.Read(buf)
	if err != nil {
		return err
	}
	b.File = buf

	return nil
}
