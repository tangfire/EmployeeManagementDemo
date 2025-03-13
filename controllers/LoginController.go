package controllers

import (
	"EmployeeManagementDemo/services"
	"EmployeeManagementDemo/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

// controllers/LoginAuth.go
type LoginRequest struct {
	Username string `json:"username" binding:"required,min=4,max=20"`
	Password string `json:"password" binding:"required,min=6,max=20"`
}

func Login(c *gin.Context) {
	var req LoginRequest

	// 参数验证
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.TranslateValidationErrors(err)})
		return
	}

	// 统一认证逻辑（同时支持管理员和员工）
	user, err := services.AuthenticateUser(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	// 生成JWT
	token, err := utils.GenerateJWT(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "令牌生成失败"})
		return
	}

	// 返回用户信息（按角色区分）
	responseData := gin.H{
		"token": token,
		"user": gin.H{
			"id":   user.GetID(),
			"name": user.GetUsername(),
			"role": user.GetRole(),
		},
	}

	c.JSON(http.StatusOK, responseData)
}
