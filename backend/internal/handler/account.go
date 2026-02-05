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

var _ = model.Account{}

type AccountRequest struct {
	Name    string `json:"name" binding:"required" example:"我的账号"`
	Session string `json:"session" binding:"required" example:"base64-session-cookie"`
}

type VerifyRequest struct {
	Session string `json:"session" binding:"required" example:"base64-session-cookie"`
}

// ListAccounts 账号列表
// @Summary 获取所有账号
// @Tags 账号管理
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=[]model.Account}
// @Router /accounts [get]
func ListAccounts(c *gin.Context) {
	accounts, err := service.ListAccounts()
	if err != nil {
		response.Error(c, 500, "获取账号失败")
		return
	}
	response.Success(c, accounts)
}

// CreateAccount 添加账号
// @Summary 添加 AnyRouter 账号
// @Tags 账号管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body AccountRequest true "账号参数"
// @Success 200 {object} response.Response{data=model.Account}
// @Router /accounts [post]
func CreateAccount(c *gin.Context) {
	var req AccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "参数错误")
		return
	}

	account, err := service.CreateAccount(req.Name, req.Session)
	if err != nil {
		if errors.Is(err, service.ErrInvalidSession) {
			response.Error(c, 400, err.Error())
			return
		}
		response.Error(c, 500, "创建失败")
		return
	}

	response.Success(c, account)
}

// UpdateAccount 更新账号
// @Summary 更新账号信息
// @Tags 账号管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "账号ID"
// @Param request body AccountRequest true "账号参数"
// @Success 200 {object} response.Response{data=model.Account}
// @Router /accounts/{id} [put]
func UpdateAccount(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, 400, "账号ID无效")
		return
	}

	var req AccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "参数错误")
		return
	}

	account, err := service.UpdateAccount(uint(id), req.Name, req.Session)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Error(c, 404, "账号不存在")
			return
		}
		if errors.Is(err, service.ErrInvalidSession) {
			response.Error(c, 400, err.Error())
			return
		}
		response.Error(c, 500, "更新失败")
		return
	}
	response.Success(c, account)
}

// DeleteAccount 删除账号
// @Summary 删除账号
// @Tags 账号管理
// @Produce json
// @Security BearerAuth
// @Param id path int true "账号ID"
// @Success 200 {object} response.Response
// @Router /accounts/{id} [delete]
func DeleteAccount(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, 400, "账号ID无效")
		return
	}
	if err := service.DeleteAccount(uint(id)); err != nil {
		response.Error(c, 500, "删除失败")
		return
	}
	response.Success(c, nil)
}

// CheckinAccount 手动签到
// @Summary 手动执行签到
// @Tags 账号管理
// @Produce json
// @Security BearerAuth
// @Param id path int true "账号ID"
// @Success 200 {object} response.Response
// @Router /accounts/{id}/checkin [post]
func CheckinAccount(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, 400, "账号ID无效")
		return
	}
	success, result := service.CheckinAccount(uint(id))
	response.Success(c, gin.H{
		"success": success,
		"result":  result,
	})
}

// VerifyAccount 验证 Session
// @Summary 验证 AnyRouter Session 有效性
// @Tags 账号管理
// @Accept json
// @Produce json
// @Param request body VerifyRequest true "Session"
// @Success 200 {object} response.Response{data=service.SessionInfo}
// @Router /accounts/verify [post]
func VerifyAccount(c *gin.Context) {
	var req VerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "参数错误")
		return
	}

	info, err := service.ParseSession(req.Session)
	if err != nil {
		response.Error(c, 400, "Session 无效: "+err.Error())
		return
	}

	response.Success(c, info)
}
