package handler

import (
	"strconv"

	"anyrouter-checkin/internal/model"
	"anyrouter-checkin/internal/repository"
	"anyrouter-checkin/internal/service"
	"anyrouter-checkin/pkg/response"

	"github.com/gin-gonic/gin"
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
	var tasks []model.CronTask
	repository.DB.Order("id desc").Find(&tasks)
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

	if err := repository.DB.Create(&task).Error; err != nil {
		response.Error(c, 500, "创建失败")
		return
	}

	service.RegisterTask(task)
	response.Success(c, task)
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
	id, _ := strconv.Atoi(c.Param("id"))
	var task model.CronTask
	if repository.DB.First(&task, id).Error != nil {
		response.Error(c, 404, "任务不存在")
		return
	}

	var req CronRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "参数错误")
		return
	}

	task.Name = req.Name
	task.CronExpr = req.CronExpr
	task.AccountIDs = req.AccountIDs
	task.Status = req.Status

	repository.DB.Save(&task)

	if task.Status == 1 {
		service.RegisterTask(task)
	} else {
		service.UnregisterTask(task.ID)
	}

	response.Success(c, task)
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
	id, _ := strconv.Atoi(c.Param("id"))
	service.UnregisterTask(uint(id))
	repository.DB.Delete(&model.CronTask{}, id)
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
	id, _ := strconv.Atoi(c.Param("id"))
	go service.ExecuteTask(uint(id))
	response.Success(c, gin.H{"message": "任务已触发"})
}
