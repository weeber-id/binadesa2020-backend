package media

import (
	"archive/zip"
	"binadesa2020-backend/lib/models"
	"binadesa2020-backend/lib/services/storage"
	"binadesa2020-backend/lib/variable"
	"bytes"
	"log"
	"net/http"
	"path"
	"sync"

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

// DownloadMultiplePrivateFile for admin
// scan files from submission search by unique code
func DownloadMultiplePrivateFile(c *gin.Context) {
	var req struct {
		UniqueCode string `form:"unique_code" binding:"required"`
	}

	// extract parameter from URL query
	if err := c.BindQuery(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	var (
		wg                                    sync.WaitGroup
		akta                                  models.AktaKelahiran
		karkel                                models.KartuKeluarga
		suratket                              models.SuratKeterangan
		foundAkta, foundKarkel, foundSuratket bool
		filesname                             []string
	)

	wg.Add(3)
	go func() {
		defer wg.Done()
		foundAkta, _ = akta.GetByUniqueCode(req.UniqueCode)
		if foundAkta == true {
			filesname = append(
				filesname,
				akta.File.KTPIstri,
				akta.File.KTPSaksi1,
				akta.File.KTPSaksi2,
				akta.File.KTPSuami,
				akta.File.SuratKelahiran,
				akta.File.SuratNikah,
			)
		}
	}()
	go func() {
		defer wg.Done()
		foundKarkel, _ = karkel.GetByUniqueCode(req.UniqueCode)
		if foundKarkel == true {
			filesname = append(
				filesname,
				karkel.File.AktaKelahiranAnak,
				karkel.File.KTPIstri,
				karkel.File.KTPSuami,
				karkel.File.SuratNikah,
			)
		}
	}()
	go func() {
		defer wg.Done()
		foundSuratket, _ = suratket.GetByUniqueCode(req.UniqueCode)
		if foundSuratket == true {
			filesname = append(
				filesname,
				suratket.File.KTP,
				suratket.File.LampiranPendukung,
				suratket.File.SuratPernyataan,
			)
		}
	}()
	wg.Wait()

	// if not found all of them
	if foundAkta == false &&
		foundKarkel == false &&
		foundSuratket == false {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "not found"})
		return
	}

	// create ZIP writer and bytes buffered for response
	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)

	// create channel
	dataChan := make(chan []byte, len(filesname))
	nameChan := make(chan string, len(filesname))

	for _, filename := range filesname {
		go func(name string) {
			var object storage.PrivateObject
			object.ObjectName = name
			data, _, err := object.Download(c)
			dataChan <- data
			nameChan <- name
			if err != nil {
				log.Printf("error in download multiple file: %s", err)
			}
		}(filename)
	}

	for i := 0; i < cap(dataChan); i++ {
		zipFile, _ := zipWriter.Create(<-nameChan)
		zipFile.Write(<-dataChan)
	}
	close(dataChan)
	close(nameChan)
	zipWriter.Close()

	contentType := http.DetectContentType(buf.Bytes())
	c.Header("Content-Disposition", "attachment; filename="+req.UniqueCode+".zip")
	c.Header("Content-Type", contentType)
	c.Data(http.StatusOK, contentType, buf.Bytes())
}
