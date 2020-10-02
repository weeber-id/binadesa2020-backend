package suratketerangan

import (
	"binadesa2020-backend/lib/clog"
	"binadesa2020-backend/lib/models"
	"binadesa2020-backend/lib/services/gmail"
	"binadesa2020-backend/lib/services/storage"
	"binadesa2020-backend/lib/tools"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin/binding"

	"github.com/gin-gonic/gin"
)

// Submission surat keterangan by user
func Submission(c *gin.Context) {
	var (
		wg  sync.WaitGroup
		req struct {
			Nama  string `form:"nama" binding:"required"`
			Email string `form:"email" binding:"required"`
			Tipe  string `form:"tipe" binding:"required"`
		}
	)

	// extract formdata string
	if err := c.ShouldBindWith(&req, binding.FormMultipart); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	// extract formdata file
	suratPernyataanHeader, err := c.FormFile("surat_pernyataan")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "surat pernyataan required"})
		return
	}
	KTPHeader, err := c.FormFile("ktp")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "KTP required"})
		return
	}
	LampiranHeader, err := c.FormFile("lampiran_pendukung")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "lampiran pendukung required"})
		return
	}

	// create new submission
	newSubmission := &models.SuratKeterangan{
		Nama:  req.Nama,
		Email: req.Email,
		Tipe:  req.Tipe,
	}
	result, err := newSubmission.Create()
	if err != nil {
		clog.Panic2Response(c, err, "create surat keterangan submission")
	}

	// object naming
	code := newSubmission.UniqueCode
	objectName := map[string]string{
		"surat_pernyataan":   "pengajuan/surat-keterangan/" + code + "/surat_pernyataan." + tools.GetExtension(suratPernyataanHeader.Filename),
		"ktp":                "pengajuan/surat-keterangan/" + code + "/ktp." + tools.GetExtension(KTPHeader.Filename),
		"lampiran_pendukung": "pengajuan/surat-keterangan/" + code + "/lampiran_pendukung." + tools.GetExtension(LampiranHeader.Filename),
	}

	// upload file to minio storage
	suratPernyataanObj := &storage.PrivateObject{}
	ktpObj := &storage.PrivateObject{}
	lampiranObj := &storage.PrivateObject{}

	wg.Add(3)
	go func() {
		defer wg.Done()
		suratPernyataanObj.LoadFileHeader(suratPernyataanHeader, objectName["surat_pernyataan"])
		suratPernyataanObj.Upload(c)
		newSubmission.File.SuratPernyataan = objectName["surat_pernyataan"]
	}()
	go func() {
		defer wg.Done()
		ktpObj.LoadFileHeader(KTPHeader, objectName["ktp"])
		ktpObj.Upload(c)
		newSubmission.File.KTP = objectName["ktp"]
	}()
	go func() {
		defer wg.Done()
		lampiranObj.LoadFileHeader(LampiranHeader, objectName["lampiran_pendukung"])
		lampiranObj.Upload(c)
		newSubmission.File.LampiranPendukung = objectName["lampiran_pendukung"]
	}()
	wg.Wait()

	_, err = newSubmission.Update()
	if err != nil {
		clog.Panic(err, "update kartu keluarga files struct")
	}

	// Sending receive email
	go func() {
		email := &gmail.Email{To: req.Email}
		email.SendReceiveSubmission("Surat Keterangan")
	}()

	c.JSON(http.StatusOK, gin.H{
		"message": "Created",
		"data": gin.H{
			"id":          result.InsertedID,
			"unique_code": newSubmission.UniqueCode,
			"is_paid":     newSubmission.IsPaid,
		},
	})
}
