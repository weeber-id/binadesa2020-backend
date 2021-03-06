package main

import (
	"binadesa2020-backend/lib/controllers"
	"binadesa2020-backend/lib/controllers/administrator"
	"binadesa2020-backend/lib/controllers/aktakelahiran"
	"binadesa2020-backend/lib/controllers/complaint"
	"binadesa2020-backend/lib/controllers/kartukeluarga"
	"binadesa2020-backend/lib/controllers/media"
	"binadesa2020-backend/lib/controllers/news"
	"binadesa2020-backend/lib/controllers/suratketerangan"
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
		root.POST("/complaint", complaint.Create)
		root.GET("/news", news.Get)

		submissionGroup := root.Group("/submission")
		{
			submissionGroup.GET("/find", controllers.GetSubmissionByCode)

			submissionGroup.GET("/kartu-keluarga", kartukeluarga.GetOne)
			submissionGroup.POST("/kartu-keluarga", kartukeluarga.Submission)

			submissionGroup.GET("/akta-kelahiran", aktakelahiran.GetOne)
			submissionGroup.POST("/akta-kelahiran", aktakelahiran.Submission)

			submissionGroup.GET("/surat-keterangan", suratketerangan.GetOne)
			submissionGroup.POST("/surat-keterangan", suratketerangan.Submission)
		}

		adminGroup := root.Group("/admin")
		adminGroup.Use(middleware.AdminAuthorization())
		{
			adminGroup.GET("/accounts", administrator.GetAll)
			adminGroup.GET("/account", administrator.GetMe)
			adminGroup.POST("/account", administrator.Create)
			adminGroup.PATCH("/account/password", administrator.ChangePassword)
			adminGroup.POST("/account/change-password", administrator.ChangePassword)
			adminGroup.DELETE("/account", administrator.Delete)
			adminGroup.POST("/account/delete", administrator.Delete)

			adminGroup.GET("/complaints", complaint.Get)
			adminGroup.POST("/complaints/read_status", complaint.ReadStatus)

			adminGroup.POST("/news", news.Create)
			adminGroup.PUT("/news", news.Update)
			adminGroup.POST("/news/update", news.Update)
			adminGroup.DELETE("/news", news.Delete)
			adminGroup.POST("/news/delete", news.Delete)

			adminGroup.POST("/media/private/download-submission", media.DownloadMultiplePrivateFile)
			adminGroup.POST("/media/private/download", media.DownloadPrivateFile)
			adminGroup.POST("/media/public/upload", media.UploadPublicFile)

			adminSubmissionGroup := adminGroup.Group("/submission")
			{
				adminSubmissionGroup.GET("/kartu-keluarga", kartukeluarga.Get)
				adminSubmissionGroup.PATCH("/kartu-keluarga", kartukeluarga.Update)
				adminSubmissionGroup.POST("/kartu-keluarga/update", kartukeluarga.Update)

				adminSubmissionGroup.GET("/akta-kelahiran", aktakelahiran.Get)
				adminSubmissionGroup.PATCH("/akta-kelahiran", aktakelahiran.Update)
				adminSubmissionGroup.POST("/akta-kelahiran/update", aktakelahiran.Update)

				adminSubmissionGroup.GET("/surat-keterangan", suratketerangan.Get)
				adminSubmissionGroup.PATCH("/surat-keterangan", suratketerangan.Update)
				adminSubmissionGroup.POST("/surat-keterangan/update", suratketerangan.Update)
			}
		}
	}

	router.Run(":8080")
}
