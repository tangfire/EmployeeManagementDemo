package controllers

import (
	"EmployeeManagementDemo/models"
	"EmployeeManagementDemo/services"
	"EmployeeManagementDemo/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AdminRegisterRequest struct {
	RegisterRequest
	SecretKey string `json:"secret_key" binding:"required"` // 必须携带密钥
}

// AdminRegister 管理员注册接口
func AdminRegister(c *gin.Context) {
	var req AdminRegisterRequest

	// 参数绑定验证
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.TranslateValidationErrors(err)})
		return
	}

	// 验证管理员密钥
	if valid := services.ValidateAdminSecret(req.SecretKey); !valid {
		c.JSON(http.StatusForbidden, gin.H{"error": "无效的管理员密钥"})
		return
	}

	// 检查用户名是否存在
	if exists, _ := services.CheckAdminNameExists(req.Username); exists {
		c.JSON(http.StatusConflict, gin.H{"error": "管理员用户名已存在"})
		return
	}

	// 密码加密
	hashedPassword, err := services.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密码处理失败"})
		return
	}

	// 构建管理员对象
	newAdmin := models.Admin{
		AdminName:     req.Username,
		AdminPassword: hashedPassword,
		AdminEmail:    req.Email,
		AdminPhone:    req.Phone,
		//Avatar:        "/avatars/default-admin.png", // 默认头像
	}

	// 创建管理员
	if err := services.CreateAdmin(&newAdmin); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建管理员失败"})
		return
	}

	// 返回创建成功响应（隐藏敏感信息）
	c.JSON(http.StatusCreated, gin.H{
		"message": "管理员账号创建成功",
		"data": gin.H{
			"admin_id": newAdmin.AdminID,
			"username": newAdmin.AdminName,
			"email":    newAdmin.AdminEmail,
		},
	})
}
