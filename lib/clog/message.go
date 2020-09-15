package clog

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Fatal log with segment name
func Fatal(err error, segmentName string) {
	if err != nil {
		log.Fatalf("Fatal in '%s' : %v", segmentName, err)
	}
}

// Panic log with segment name
func Panic(err error, segmentName string) {
	if err != nil {
		log.Panicf("Panic in '%s' : %v", segmentName, err)
	}
}

// Panic2Response with segment name
func Panic2Response(c *gin.Context, err error, segmentName string) {
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		log.Panicf("Panic in '%s' : %v", segmentName, err)
	}
}
