package handler

import (
	"errors"

	"anyrouter-checkin/internal/model"
	"anyrouter-checkin/internal/repository"
	"anyrouter-checkin/internal/service"
	"anyrouter-checkin/pkg/response"

	"github.com/gin-gonic/gin"
)

// GetConfigs 获取配置
// @Summary 获取指定分类的配置
// @Tags 系统配置
// @Produce json
// @Security BearerAuth
// @Param category path string true "配置分类" Enums(telegram, system)
// @Success 200 {object} response.Response{data=map[string]string}
// @Router /config/{category} [get]
func GetConfigs(c *gin.Context) {
	category := c.Param("category")
	var configs []model.Config
	repository.DB.Where("category = ?", category).Find(&configs)

	result := make(map[string]string)
	for _, cfg := range configs {
		result[cfg.Key] = cfg.Value
	}
	response.Success(c, result)
}

// UpdateConfigs 更新配置
// @Summary 更新指定分类的配置
// @Tags 系统配置
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param category path string true "配置分类" Enums(telegram, system)
// @Param request body map[string]string true "配置键值对"
// @Success 200 {object} response.Response
// @Router /config/{category} [put]
func UpdateConfigs(c *gin.Context) {
	category := c.Param("category")
	var req map[string]string
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "参数错误")
		return
	}

	for key, value := range req {
		service.SetConfig(key, value, category)
	}

	response.Success(c, nil)
}

// TestTelegram 测试推送
// @Summary 发送 Telegram 测试消息（使用最近成功签到记录）
// @Tags 系统配置
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response
// @Router /config/telegram/test [post]
func TestTelegram(c *gin.Context) {
	err := service.SendTestCheckinNotification()
	if err != nil {
		if errors.Is(err, service.ErrNoSuccessfulCheckinLog) {
			response.Error(c, 400, err.Error())
			return
		}
		response.Error(c, 500, "发送失败: "+err.Error())
		return
	}
	response.Success(c, gin.H{"message": "发送成功"})
}

// ListLogs 签到日志
// @Summary 获取签到日志列表与今日账号统计
// @Tags 日志
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=service.CheckinLogSummary}
// @Router /logs [get]
func ListLogs(c *gin.Context) {
	data, err := service.GetCheckinLogSummary(100)
	if err != nil {
		response.Error(c, 500, "获取签到日志失败")
		return
	}
	response.Success(c, data)
}
