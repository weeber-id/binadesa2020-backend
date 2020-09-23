package kartukeluarga

import (
	"binadesa2020-backend/lib/clog"
	"binadesa2020-backend/lib/models"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin/binding"

	"github.com/gin-gonic/gin"
)

// Submission kartu keluarga
func Submission(c *gin.Context) {
	var req struct {
		NamaKepalaKeluarga string `form:"nama_kepala_keluarga" binding:"required"`
		NoTelp             string `form:"no_telp" binding:"required"`
		// KTPSuami           string `form:"ktp_suami"`
		// KTPIstri           string `form:"ktp_istri"`
		// SuratNikah         string `form:"surat_nikah"`
		// AktaKelahiranAnak  string `form:"akta_kelahiran_anak"`
	}
	log.Println("start")

	// Extract formdata string
	if err := c.ShouldBindWith(&req, binding.FormMultipart); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	// Extract formdata file
	ktpSuami, err := c.FormFile("ktp_suami")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "required ktp_suami"})
		return
	}
	ktpIstri, err := c.FormFile("ktp_istri")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "required ktp_istri"})
		return
	}
	suratNikah, err := c.FormFile("surat_nikah")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "required surat_nikah"})
		return
	}

	fmt.Printf("file ktp %v \n", ktpSuami.Filename)
	fmt.Printf("file ktp %v \n", ktpIstri.Filename)
	fmt.Printf("file surat_nikah %v \n", suratNikah.Filename)

	newSubmission := &models.KartuKeluarga{
		NamaKepalaKeluarga: req.NamaKepalaKeluarga,
		NoTelp:             req.NoTelp,
	}
	result, err := newSubmission.Create()
	if err != nil {
		clog.Panic(err, "create kartukeluarga submission")
	}

	newFiles := &models.KartuKeluargaFiles{
		KTPSuami:   ktpSuami.Filename,
		KTPIstri:   ktpIstri.Filename,
		SuratNikah: suratNikah.Filename,
	}
	_, err = newSubmission.ChangeAllFiles(newFiles)
	if err != nil {
		clog.Panic(err, "change kartukeluarga all files")
	}

	c.JSON(http.StatusOK, gin.H{"message": "Created", "data": result})
}
