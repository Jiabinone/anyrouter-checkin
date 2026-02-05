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

	r.Any("/anyrouter/*path", handler.AnyRouterProxy)

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
			auth.PUT("/accounts/:id/status", handler.UpdateAccountStatus)
			auth.DELETE("/accounts/:id", handler.DeleteAccount)
			auth.POST("/accounts/:id/checkin", handler.CheckinAccount)
			auth.POST("/accounts/:id/refresh", handler.RefreshAccount)

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
