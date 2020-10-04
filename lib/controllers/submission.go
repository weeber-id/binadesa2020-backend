package controllers

import (
	"binadesa2020-backend/lib/models"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

// GetSubmissionByCode for user
func GetSubmissionByCode(c *gin.Context) {
	var (
		wg  sync.WaitGroup
		req struct {
			UniqueCode string `form:"unique_code" binding:"required"`
		}
	)

	// extract query parameter
	if err := c.BindQuery(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	var (
		karkel models.KartuKeluarga
		akta   models.AktaKelahiran
		sK     models.SuratKeterangan
		data   struct {
			AktaKelahiran   *models.AktaKelahiran   `json:"akta_kelahiran,omitempty"`
			KartuKeluarga   *models.KartuKeluarga   `json:"kartu_keluarga,omitempty"`
			SuratKeterangan *models.SuratKeterangan `json:"surat_keterangan,omitempty"`
		}
	)

	// Search by unique code concurrently
	wg.Add(3)
	go func() {
		defer wg.Done()
		akta.GetByUniqueCode(req.UniqueCode)
		if (akta != models.AktaKelahiran{}) {
			data.AktaKelahiran = &akta
		}
	}()
	go func() {
		defer wg.Done()
		karkel.GetByUniqueCode(req.UniqueCode)
		if (karkel != models.KartuKeluarga{}) {
			data.KartuKeluarga = &karkel
		}
	}()
	go func() {
		defer wg.Done()
		sK.GetByUniqueCode(req.UniqueCode)
		if (sK != models.SuratKeterangan{}) {
			data.SuratKeterangan = &sK
		}
	}()
	wg.Wait()

	// if not found all of them
	if data.AktaKelahiran == nil &&
		data.KartuKeluarga == nil &&
		data.SuratKeterangan == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "not found"})
		return
	}
	// if found
	c.JSON(http.StatusOK, gin.H{"message": "OK", "data": data})
}
