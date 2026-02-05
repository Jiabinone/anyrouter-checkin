package handler

import (
	"errors"
	"io"
	"strconv"

	"anyrouter-checkin/internal/service"
	"anyrouter-checkin/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CreateAccountRequest struct {
	Session string `json:"session" binding:"required" example:"base64-session-cookie"`
}

type UpdateAccountRequest struct {
	Session string `json:"session" example:"base64-session-cookie"`
}

type UpdateAccountStatusRequest struct {
	Status *int `json:"status" binding:"required" example:"1"`
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
// @Param request body CreateAccountRequest true "账号参数"
// @Success 200 {object} response.Response{data=model.Account}
// @Router /accounts [post]
func CreateAccount(c *gin.Context) {
	var req CreateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "参数错误")
		return
	}

	account, err := service.CreateAccount(req.Session)
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
// @Param request body UpdateAccountRequest true "账号参数"
// @Success 200 {object} response.Response{data=model.Account}
// @Router /accounts/{id} [put]
func UpdateAccount(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, 400, "账号ID无效")
		return
	}

	var req UpdateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil && !errors.Is(err, io.EOF) {
		response.Error(c, 400, "参数错误")
		return
	}

	account, err := service.UpdateAccount(uint(id), req.Session)
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

// UpdateAccountStatus 更新账号状态
// @Summary 更新账号状态
// @Tags 账号管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "账号ID"
// @Param request body UpdateAccountStatusRequest true "账号状态"
// @Success 200 {object} response.Response{data=model.Account}
// @Router /accounts/{id}/status [put]
func UpdateAccountStatus(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, 400, "账号ID无效")
		return
	}

	var req UpdateAccountStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "参数错误")
		return
	}
	if req.Status == nil || (*req.Status != 0 && *req.Status != 1) {
		response.Error(c, 400, "状态无效")
		return
	}

	account, err := service.UpdateAccountStatus(uint(id), *req.Status)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Error(c, 404, "账号不存在")
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

// RefreshAccount 刷新账号信息
// @Summary 刷新账号信息
// @Tags 账号管理
// @Produce json
// @Security BearerAuth
// @Param id path int true "账号ID"
// @Success 200 {object} response.Response{data=model.Account}
// @Router /accounts/{id}/refresh [post]
func RefreshAccount(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, 400, "账号ID无效")
		return
	}

	account, err := service.RefreshAccount(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Error(c, 404, "账号不存在")
			return
		}
		if errors.Is(err, service.ErrInvalidSession) {
			response.Error(c, 400, err.Error())
			return
		}
		if errors.Is(err, service.ErrAccountDisabled) {
			response.Error(c, 400, err.Error())
			return
		}
		response.Error(c, 500, "刷新失败")
		return
	}

	response.Success(c, account)
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
