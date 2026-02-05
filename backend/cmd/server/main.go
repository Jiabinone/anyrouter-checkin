package main

import (
	"fmt"
	"os"
	"path/filepath"

	"anyrouter-checkin/internal/config"
	"anyrouter-checkin/internal/repository"
	"anyrouter-checkin/internal/router"
	"anyrouter-checkin/internal/service"
	"anyrouter-checkin/pkg/logger"

	_ "anyrouter-checkin/docs"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// @title AnyRouter 管理后台 API
// @version 1.0
// @description AnyRouter 签到中台管理系统
// @host localhost:8080
// @BasePath /api
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	if err := config.Load(); err != nil {
		fallback, _ := zap.NewDevelopment()
		fallback.Fatal("加载配置失败", zap.Error(err))
	}

	zapLogger, err := logger.Init(config.C.Server.Mode)
	if err != nil {
		fallback, _ := zap.NewDevelopment()
		fallback.Fatal("初始化日志失败", zap.Error(err))
	}
	defer func() {
		_ = zapLogger.Sync()
	}()

	dbDir := filepath.Dir(config.C.Database.Path)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		zap.L().Fatal("创建数据目录失败", zap.Error(err))
	}

	if err := repository.Init(config.C.Database.Path); err != nil {
		zap.L().Fatal("初始化数据库失败", zap.Error(err))
	}

	repository.InitDefaultConfigs()
	if err := service.InitAdminUser(); err != nil {
		zap.L().Fatal("初始化管理员失败", zap.Error(err))
	}

	service.InitCron()

	gin.SetMode(config.C.Server.Mode)
	gin.DefaultWriter = logger.Writer(zapLogger, zapcore.InfoLevel)
	gin.DefaultErrorWriter = logger.Writer(zapLogger, zapcore.ErrorLevel)
	r := gin.Default()
	if err := r.SetTrustedProxies(nil); err != nil {
		zap.L().Fatal("设置 TrustedProxies 失败", zap.Error(err))
	}
	router.Setup(r)

	addr := fmt.Sprintf(":%d", config.C.Server.Port)
	zap.L().Info("服务启动", zap.String("addr", "http://localhost"+addr))
	zap.L().Info("Swagger", zap.String("url", "http://localhost"+addr+"/swagger/index.html"))
	if err := r.Run(addr); err != nil {
		zap.L().Fatal("启动失败", zap.Error(err))
	}
}
