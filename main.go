package main

import (
	"binadesa2020-backend/lib/controllers"
	"binadesa2020-backend/lib/controllers/kartukeluarga"
	"binadesa2020-backend/lib/middleware"
	"binadesa2020-backend/lib/services/mongodb"
	"binadesa2020-backend/lib/services/storage"
	"binadesa2020-backend/lib/variable"
	"time"

	"context"

	"github.com/gin-gonic/gin"
)

func main() {
	variable.Initialization()
	storage.MinioInitialization()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	mongodb.Connection(ctx)
	defer mongodb.Client.Disconnect(ctx)

	router := gin.Default()
	router.POST("/login", controllers.Login)

	adminGroup := router.Group("/admin")
	adminGroup.Use(middleware.AdminAuthorization())
	{
		adminGroup.GET("/accounts", controllers.GetAllAdmin)
		adminGroup.POST("/account", controllers.CreateAdmin)
		adminGroup.DELETE("/account", controllers.DeleteAdmin)

		adminGroup.GET("/complaints", controllers.GetAllComplaint)
	}

	router.POST("/complaint", controllers.CreateComplaint)
	submissionGroup := router.Group("/submission")
	{
		submissionGroup.POST("/kartu-keluarga", kartukeluarga.Submission)
	}

	router.Run(":8080")
}
