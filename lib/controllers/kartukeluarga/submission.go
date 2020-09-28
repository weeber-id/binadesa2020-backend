package kartukeluarga

import (
	"binadesa2020-backend/lib/clog"
	"binadesa2020-backend/lib/models"
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
		"ktp_suami":   "pengajuan/kartu-keluarga/" + code + "/ktpsuami." + tools.GetExtension(ktpSuamiHeader.Filename),
		"ktp_istri":   "pengajuan/kartu-keluarga/" + code + "/ktpistri." + tools.GetExtension(ktpIstriHeader.Filename),
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
	}()
	go func() {
		defer wg.Done()
		ktpIstriObj.LoadFileHeader(ktpIstriHeader, objectName["ktp_istri"])
		ktpIstriObj.Upload(c)
	}()
	go func() {
		defer wg.Done()
		suratNikahObj.LoadFileHeader(suratNikahHeader, objectName["surat_nikah"])
		suratNikahObj.Upload(c)
	}()
	go func() {
		defer wg.Done()
		aktaKelahiranAnakObj.LoadFileHeader(aktaKelahiranAnakHeader, objectName["akta_kelahiran_anak"])
		aktaKelahiranAnakObj.Upload(c)
	}()
	wg.Wait()

	newFiles := &models.KartuKeluargaFiles{
		KTPSuami:          objectName["ktp_suami"],
		KTPIstri:          objectName["ktp_istri"],
		SuratNikah:        objectName["surat_nikah"],
		AktaKelahiranAnak: objectName["akta_kelahiran_anak"],
	}
	_, err = newSubmission.ChangeAllFiles(newFiles)
	if err != nil {
		clog.Panic(err, "change kartukeluarga all files")
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Created",
		"data": gin.H{
			"id":          result.InsertedID,
			"unique_code": newSubmission.UniqueCode,
		},
	})
}
