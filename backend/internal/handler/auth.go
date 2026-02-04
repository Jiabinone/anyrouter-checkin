package handler

import (
	"anyrouter-checkin/internal/service"
	"anyrouter-checkin/pkg/response"

	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required" example:"admin"`
	Password string `json:"password" binding:"required" example:"admin123"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required" example:"oldpass"`
	NewPassword string `json:"new_password" binding:"required,min=6" example:"newpass123"`
}

// Login 用户登录
// @Summary 用户登录
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body LoginRequest true "登录参数"
// @Success 200 {object} response.Response{data=map[string]string}
// @Router /auth/login [post]
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "参数错误")
		return
	}

	token, err := service.Login(req.Username, req.Password)
	if err != nil {
		response.Error(c, 401, err.Error())
		return
	}

	response.Success(c, gin.H{"token": token})
}

// Profile 获取当前用户信息
// @Summary 获取当前用户信息
// @Tags 认证
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response
// @Router /auth/profile [get]
func Profile(c *gin.Context) {
	userID, _ := c.Get("user_id")
	username, _ := c.Get("username")
	response.Success(c, gin.H{
		"user_id":  userID,
		"username": username,
	})
}

// ChangePassword 修改密码
// @Summary 修改当前用户密码
// @Tags 认证
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body ChangePasswordRequest true "密码参数"
// @Success 200 {object} response.Response
// @Router /auth/password [put]
func ChangePassword(c *gin.Context) {
	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "参数错误：新密码至少6位")
		return
	}

	userID, _ := c.Get("user_id")
	if err := service.ChangePassword(userID.(uint), req.OldPassword, req.NewPassword); err != nil {
		response.Error(c, 400, err.Error())
		return
	}

	response.Success(c, nil)
}
