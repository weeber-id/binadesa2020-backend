package main

import (
	"binadesa2020-backend/lib/controllers"
	"binadesa2020-backend/lib/controllers/aktakelahiran"
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
	root := router.Group("/api")
	root.Use(middleware.CORS())
	{
		root.GET("/", controllers.HealthCheck)
		root.POST("/login", controllers.Login)
		root.POST("/logout", controllers.Logout)
		root.POST("/complaint", controllers.CreateComplaint)

		adminGroup := root.Group("/admin")
		adminGroup.Use(middleware.AdminAuthorization())
		{
			adminGroup.GET("/accounts", controllers.GetAllAdmin)
			adminGroup.POST("/account", controllers.CreateAdmin)
			adminGroup.DELETE("/account", controllers.DeleteAdmin)

			adminGroup.GET("/complaints", controllers.GetAllComplaint)

			adminSubmissionGroup := adminGroup.Group("/submission")
			{
				adminSubmissionGroup.GET("/kartu-keluarga", kartukeluarga.Get)
				adminSubmissionGroup.PATCH("/kartu-keluarga", kartukeluarga.Update)

				adminSubmissionGroup.GET("/akta-kelahiran", aktakelahiran.Get)
				adminSubmissionGroup.PATCH("/akta-kelahiran", aktakelahiran.Update)
			}
		}

		submissionGroup := root.Group("/submission")
		{
			submissionGroup.GET("/kartu-keluarga", kartukeluarga.GetOne)
			submissionGroup.POST("/kartu-keluarga", kartukeluarga.Submission)

			submissionGroup.GET("/akta-kelahiran", aktakelahiran.GetOne)
			submissionGroup.POST("/akta-kelahiran", aktakelahiran.Submission)
		}
	}

	router.Run(":8080")
}
