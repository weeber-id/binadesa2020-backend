package kartukeluarga

import (
	"binadesa2020-backend/lib/clog"
	"binadesa2020-backend/lib/models"
	"binadesa2020-backend/lib/services/gmail"
	"binadesa2020-backend/lib/services/storage"
	"binadesa2020-backend/lib/tools"
	"mime/multipart"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin/binding"

	"github.com/gin-gonic/gin"
)

// Submission kartu keluarga
func Submission(c *gin.Context) {
	var (
		wg  sync.WaitGroup
		req struct {
			Nama               string `form:"nama" binding:"required"`
			NamaKepalaKeluarga string `form:"nama_kepala_keluarga" binding:"required"`
			Email              string `form:"email" binding:"required"`
		}
	)

	// Extract formdata string
	if err := c.ShouldBindWith(&req, binding.FormMultipart); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	// Extract formdata file
	ktpSuamiHeader, err := c.FormFile("ktp_suami")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "required ktp_suami"})
		return
	}
	ktpIstriHeader, err := c.FormFile("ktp_istri")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "required ktp_istri"})
		return
	}
	suratNikahHeader, err := c.FormFile("surat_nikah")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "required surat_nikah"})
		return
	}
	aktaKelahiranAnakHeader, _ := c.FormFile("akta_kelahiran_anak")
	if aktaKelahiranAnakHeader == nil {
		aktaKelahiranAnakHeader = &multipart.FileHeader{}
	}

	// create new submission
	newSubmission := &models.KartuKeluarga{
		Nama:               req.Nama,
		NamaKepalaKeluarga: req.NamaKepalaKeluarga,
		Email:              req.Email,
	}
	result, err := newSubmission.Create()
	if err != nil {
		clog.Panic(err, "create kartukeluarga submission")
	}

	// object naming
	code := newSubmission.UniqueCode
	objectName := map[string]string{
		"ktp_suami":   "pengajuan/kartu-keluarga/" + code + "/ktp_suami." + tools.GetExtension(ktpSuamiHeader.Filename),
		"ktp_istri":   "pengajuan/kartu-keluarga/" + code + "/ktp_istri." + tools.GetExtension(ktpIstriHeader.Filename),
		"surat_nikah": "pengajuan/kartu-keluarga/" + code + "/surat_nikah." + tools.GetExtension(suratNikahHeader.Filename),
	}
	if aktaKelahiranAnakHeader.Filename != "" {
		objectName["akta_kelahiran_anak"] = "pengajuan/kartu-keluarga/" + code + "/akta_kelahiran_anak." + tools.GetExtension(aktaKelahiranAnakHeader.Filename)
	} else {
		objectName["akta_kelahiran_anak"] = ""
	}

	// upload file to minio storage
	ktpSuamiObj := &storage.PrivateObject{}
	ktpIstriObj := &storage.PrivateObject{}
	suratNikahObj := &storage.PrivateObject{}
	aktaKelahiranAnakObj := &storage.PrivateObject{}

	wg.Add(4)
	go func() {
		defer wg.Done()
		ktpSuamiObj.LoadFileHeader(ktpSuamiHeader, objectName["ktp_suami"])
		ktpSuamiObj.Upload(c)
		newSubmission.File.KTPSuami = objectName["ktp_suami"]
	}()
	go func() {
		defer wg.Done()
		ktpIstriObj.LoadFileHeader(ktpIstriHeader, objectName["ktp_istri"])
		ktpIstriObj.Upload(c)
		newSubmission.File.KTPIstri = objectName["ktp_istri"]
	}()
	go func() {
		defer wg.Done()
		suratNikahObj.LoadFileHeader(suratNikahHeader, objectName["surat_nikah"])
		suratNikahObj.Upload(c)
		newSubmission.File.SuratNikah = objectName["surat_nikah"]
	}()
	go func() {
		defer wg.Done()
		aktaKelahiranAnakObj.LoadFileHeader(aktaKelahiranAnakHeader, objectName["akta_kelahiran_anak"])
		aktaKelahiranAnakObj.Upload(c)
		newSubmission.File.AktaKelahiranAnak = objectName["akta_kelahiran_anak"]
	}()
	wg.Wait()

	_, err = newSubmission.Update()
	if err != nil {
		clog.Panic(err, "update kartu keluarga files struct")
	}

	// Sending receive email
	go func() {
		email := &gmail.Email{To: req.Email}
		email.SendReceiveSubmission(newSubmission)
	}()

	c.JSON(http.StatusOK, gin.H{
		"message": "Created",
		"data": gin.H{
			"id":          result.InsertedID,
			"unique_code": newSubmission.UniqueCode,
		},
	})
}
