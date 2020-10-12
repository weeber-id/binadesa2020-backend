package complaint

import (
	"binadesa2020-backend/lib/models"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

func ReadStatus(c *gin.Context) {
	var (
		wg  sync.WaitGroup
		req struct {
			IDs    []string `json:"ids" binding:"required"`
			IsRead *bool    `json:"is_read" binding:"required"`
		}
	)

	// extract parameter
	if err := c.BindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	var status struct {
		IsRead      bool `json:"is_read"`
		NumFound    int  `json:"num_found"`
		NumNotFound int  `json:"num_not_found"`
	}

	wg.Add(len(req.IDs))
	for _, id := range req.IDs {
		go func(ID string) {
			defer wg.Done()
			var complaint models.Complaint
			found, _ := complaint.GetByID(ID)

			if found == false {
				status.NumNotFound++
				return
			}

			complaint.IsRead = *req.IsRead
			complaint.Update()
			status.NumFound++
		}(id)
	}
	wg.Wait()

	c.JSON(http.StatusOK, gin.H{"message": "OK", "data": status})
}
