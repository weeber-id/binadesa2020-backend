package main

import (
	"binadesa2020-backend/lib/controllers"
	"binadesa2020-backend/lib/services/mongodb"
	"binadesa2020-backend/lib/variable"
	"time"

	"context"

	"github.com/gin-gonic/gin"
)

func main() {
	variable.Initialization()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	mongodb.Connection(ctx)
	defer mongodb.Client.Disconnect(ctx)

	router := gin.Default()
	adminGroup := router.Group("/admin")
	{
		adminGroup.POST("/login", controllers.Login)
		adminGroup.GET("/list", controllers.GetAdmin)
	}

	router.Run(":8080")
}
