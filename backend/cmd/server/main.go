package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"anyrouter-checkin/internal/config"
	"anyrouter-checkin/internal/repository"
	"anyrouter-checkin/internal/router"
	"anyrouter-checkin/internal/service"

	_ "anyrouter-checkin/docs"

	"github.com/gin-gonic/gin"
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
		log.Fatalf("加载配置失败: %v", err)
	}

	dbDir := filepath.Dir(config.C.Database.Path)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		log.Fatalf("创建数据目录失败: %v", err)
	}

	if err := repository.Init(config.C.Database.Path); err != nil {
		log.Fatalf("初始化数据库失败: %v", err)
	}

	repository.InitDefaultConfigs()
	if err := service.InitAdminUser(); err != nil {
		log.Fatalf("初始化管理员失败: %v", err)
	}

	service.InitCron()

	gin.SetMode(config.C.Server.Mode)
	r := gin.Default()
	router.Setup(r)

	addr := fmt.Sprintf(":%d", config.C.Server.Port)
	log.Printf("服务启动: http://localhost%s", addr)
	log.Printf("Swagger: http://localhost%s/swagger/index.html", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("启动失败: %v", err)
	}
}
