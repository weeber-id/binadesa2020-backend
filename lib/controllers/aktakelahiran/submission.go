package aktakelahiran

import (
	"binadesa2020-backend/lib/clog"
	"binadesa2020-backend/lib/models"
	"binadesa2020-backend/lib/services/storage"
	"binadesa2020-backend/lib/tools"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// Submission aktakelahiran by user
func Submission(c *gin.Context) {
	var (
		wg  sync.WaitGroup
		req struct {
			NamaKepalaKeluarga string `form:"nama_kepala_keluarga" binding:"required"`
			Email              string `form:"email" binding:"required"`
		}
	)

	// extract formdata string
	if err := c.ShouldBindWith(&req, binding.FormMultipart); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	// extract formdata file
	suratKelahiranHeader, err := c.FormFile("surat_kelahiran")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "required surat kelahiran"})
		return
	}
	ktpSuamiHeader, err := c.FormFile("ktp_suami")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "required ktp suami"})
		return
	}
	ktpIstriHeader, err := c.FormFile("ktp_istri")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "required ktp istri"})
		return
	}
	ktpSaksi1Header, err := c.FormFile("ktp_saksi_1")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "required ktp saksi 1"})
		return
	}
	ktpSaksi2Header, err := c.FormFile("ktp_saksi_2")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "required ktp saksi 2"})
		return
	}
	suratNikahHeader, err := c.FormFile("surat_nikah")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "required surat nikah"})
		return
	}

	// craete new submission
	newSubmission := &models.AktaKelahiran{
		NamaKepalaKeluarga: req.NamaKepalaKeluarga,
		Email:              req.Email,
	}
	result, err := newSubmission.Create()
	if err != nil {
		clog.Panic(err, "create akta kelahiran")
	}

	// object naming
	code := newSubmission.UniqueCode
	objectName := map[string]string{
		"surat_kelahiran": "pengajuan/akta-kelahiran/" + code + "/surat_kelahiran." + tools.GetExtension(suratKelahiranHeader.Filename),
		"ktp_suami":       "pengajuan/akta-kelahiran/" + code + "/ktp_suami." + tools.GetExtension(ktpSuamiHeader.Filename),
		"ktp_istri":       "pengajuan/akta-kelahiran/" + code + "/ktp_istri." + tools.GetExtension(ktpIstriHeader.Filename),
		"ktp_saksi_1":     "pengajuan/akta-kelahiran/" + code + "/ktp_saksi_1." + tools.GetExtension(ktpSaksi1Header.Filename),
		"ktp_saksi_2":     "pengajuan/akta-kelahiran/" + code + "/ktp_saksi_2." + tools.GetExtension(ktpSaksi2Header.Filename),
		"surat_nikah":     "pengajuan/akta-kelahiran/" + code + "/surat_nikah." + tools.GetExtension(suratNikahHeader.Filename),
	}

	// Upload files to minio storage
	suratKelahiranObj := &storage.PrivateObject{}
	ktpSuamiObj := &storage.PrivateObject{}
	ktpIstriObj := &storage.PrivateObject{}
	ktpSaksi1Obj := &storage.PrivateObject{}
	ktpSaksi2Obj := &storage.PrivateObject{}
	suratNikahObj := &storage.PrivateObject{}

	wg.Add(6)
	go func() {
		defer wg.Done()
		suratKelahiranObj.LoadFileHeader(suratKelahiranHeader, objectName["surat_kelahiran"])
		suratKelahiranObj.Upload(c)
		newSubmission.File.SuratKelahiran = objectName["surat_kelahiran"]
	}()
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
		ktpSaksi1Obj.LoadFileHeader(ktpSaksi1Header, objectName["ktp_saksi_1"])
		ktpSaksi1Obj.Upload(c)
		newSubmission.File.KTPSaksi1 = objectName["ktp_saksi_1"]
	}()
	go func() {
		defer wg.Done()
		ktpSaksi2Obj.LoadFileHeader(ktpSaksi2Header, objectName["ktp_saksi_2"])
		ktpSaksi2Obj.Upload(c)
		newSubmission.File.KTPSaksi2 = objectName["ktp_saksi_2"]
	}()
	go func() {
		defer wg.Done()
		suratNikahObj.LoadFileHeader(suratNikahHeader, objectName["surat_nikah"])
		suratNikahObj.Upload(c)
		newSubmission.File.SuratNikah = objectName["surat_nikah"]
	}()
	wg.Wait()

	_, err = newSubmission.Update()
	if err != nil {
		clog.Panic(err, "change aktakelahiran struct")
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Created",
		"data": gin.H{
			"id":          result.InsertedID,
			"unique_code": newSubmission.UniqueCode,
		},
	})
}
