package router

import (
	"anyrouter-checkin/internal/handler"
	"anyrouter-checkin/internal/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Setup(r *gin.Engine) {
	r.Use(middleware.CORS())

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api")
	{
		api.POST("/auth/login", handler.Login)
		api.POST("/accounts/verify", handler.VerifyAccount)

		auth := api.Group("")
		auth.Use(middleware.Auth())
		{
			auth.GET("/auth/profile", handler.Profile)
			auth.PUT("/auth/password", handler.ChangePassword)

			auth.GET("/accounts", handler.ListAccounts)
			auth.POST("/accounts", handler.CreateAccount)
			auth.PUT("/accounts/:id", handler.UpdateAccount)
			auth.DELETE("/accounts/:id", handler.DeleteAccount)
			auth.POST("/accounts/:id/checkin", handler.CheckinAccount)

			auth.GET("/cron", handler.ListCronTasks)
			auth.POST("/cron", handler.CreateCronTask)
			auth.PUT("/cron/:id", handler.UpdateCronTask)
			auth.DELETE("/cron/:id", handler.DeleteCronTask)
			auth.POST("/cron/:id/trigger", handler.TriggerCronTask)

			auth.GET("/config/:category", handler.GetConfigs)
			auth.PUT("/config/:category", handler.UpdateConfigs)
			auth.POST("/config/telegram/test", handler.TestTelegram)

			auth.GET("/logs", handler.ListLogs)
		}
	}
}
