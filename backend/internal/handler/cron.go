package handler

import (
	"errors"
	"strconv"

	"anyrouter-checkin/internal/model"
	"anyrouter-checkin/internal/service"
	"anyrouter-checkin/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CronRequest struct {
	Name       string `json:"name" binding:"required" example:"每日签到"`
	CronExpr   string `json:"cron_expr" binding:"required" example:"0 8 * * *"`
	TaskType   string `json:"task_type" example:"checkin"`
	AccountIDs string `json:"account_ids" example:"[1,2]"`
	Status     int    `json:"status" example:"1"`
}

// ListCronTasks 定时任务列表
// @Summary 获取所有定时任务
// @Tags 定时任务
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=[]model.CronTask}
// @Router /cron [get]
func ListCronTasks(c *gin.Context) {
	tasks, err := service.ListCronTasks()
	if err != nil {
		response.Error(c, 500, "获取任务失败")
		return
	}
	response.Success(c, tasks)
}

// CreateCronTask 创建定时任务
// @Summary 创建定时任务
// @Tags 定时任务
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CronRequest true "任务参数"
// @Success 200 {object} response.Response{data=model.CronTask}
// @Router /cron [post]
func CreateCronTask(c *gin.Context) {
	var req CronRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "参数错误")
		return
	}

	task := model.CronTask{
		Name:       req.Name,
		CronExpr:   req.CronExpr,
		TaskType:   "checkin",
		AccountIDs: req.AccountIDs,
		Status:     1,
	}

	created, err := service.CreateCronTask(task)
	if err != nil {
		response.Error(c, 500, "创建失败")
		return
	}

	response.Success(c, created)
}

// UpdateCronTask 更新定时任务
// @Summary 更新定时任务
// @Tags 定时任务
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "任务ID"
// @Param request body CronRequest true "任务参数"
// @Success 200 {object} response.Response{data=model.CronTask}
// @Router /cron/{id} [put]
func UpdateCronTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, 400, "任务ID无效")
		return
	}

	var req CronRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "参数错误")
		return
	}

	updated, err := service.UpdateCronTask(uint(id), model.CronTask{
		Name:       req.Name,
		CronExpr:   req.CronExpr,
		TaskType:   "checkin",
		AccountIDs: req.AccountIDs,
		Status:     req.Status,
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Error(c, 404, "任务不存在")
			return
		}
		response.Error(c, 500, "更新失败")
		return
	}
	response.Success(c, updated)
}

// DeleteCronTask 删除定时任务
// @Summary 删除定时任务
// @Tags 定时任务
// @Produce json
// @Security BearerAuth
// @Param id path int true "任务ID"
// @Success 200 {object} response.Response
// @Router /cron/{id} [delete]
func DeleteCronTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, 400, "任务ID无效")
		return
	}
	if err := service.DeleteCronTask(uint(id)); err != nil {
		response.Error(c, 500, "删除失败")
		return
	}
	response.Success(c, nil)
}

// TriggerCronTask 立即执行任务
// @Summary 立即触发执行定时任务
// @Tags 定时任务
// @Produce json
// @Security BearerAuth
// @Param id path int true "任务ID"
// @Success 200 {object} response.Response
// @Router /cron/{id}/trigger [post]
func TriggerCronTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, 400, "任务ID无效")
		return
	}
	go service.ExecuteTask(uint(id))
	response.Success(c, gin.H{"message": "任务已触发"})
}
