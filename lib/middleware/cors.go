package middleware

import (
	"log"

	"github.com/gin-gonic/gin"
)

// CORS for website
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("Request Host: %s", c.Request.Host)
		switch c.Request.Host {
		case "localhost:3000", "staging-binadesa.weeber.id", "telukjambe.id":
			c.Writer.Header().Set("Access-Control-Allow-Origin", c.Request.Host)
		}

		c.Writer.Header().Set("Content-Type", "application/json")
		// c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")
		c.Writer.Header().Set("Access-Control-Max-Age", "600")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
