package controllers

import (
	"EmployeeManagementDemo/config"
	"EmployeeManagementDemo/models"
	"EmployeeManagementDemo/services"
	"EmployeeManagementDemo/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
	"strings"
	"time"
)

func Login(c *gin.Context) {
	var req models.LoginRequest

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

	// 发送登录日志消息
	logData := map[string]interface{}{
		"user_id":   user.GetID(),
		"action":    "login",
		"target_id": "", // 登录操作无目标对象
	}
	services.SendLogToRabbitMQ(logData)

	c.JSON(http.StatusOK, responseData)
}

// controllers/user.go
func GetProfile(c *gin.Context) {
	// 从上下文中获取用户身份
	userID, _ := utils.GetCurrentUserID(c)
	role, _ := utils.GetCurrentUserRole(c)

	var profile interface{}
	var err error

	// 根据角色查询不同表
	if role == "admin" {
		var admin models.Admin
		err = config.DB.First(&admin, userID).Error
		profile = admin.ToProfileResponse()
	} else {
		var emp models.Employee
		err = config.DB.First(&emp, userID).Error
		profile = emp.ToProfileResponse()
	}

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	c.JSON(http.StatusOK, profile)
}

// controllers/user.go
func UpdateProfile(c *gin.Context) {
	// 从上下文中获取用户身份
	userID, _ := utils.GetCurrentUserID(c)
	role, _ := utils.GetCurrentUserRole(c)

	// 绑定请求参数
	var req models.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误: " + err.Error()})
		return
	}

	// 根据角色获取用户模型
	var err error
	if role == "admin" {
		var admin models.Admin
		if err = config.DB.First(&admin, userID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "管理员不存在"})
			return
		}
		// 更新允许修改的字段 (避免直接覆盖敏感字段如密码)
		admin.AdminName = req.Name
		admin.AdminEmail = req.Email
		admin.AdminPhone = req.Phone
		admin.Avatar = req.Avatar
		err = config.DB.Save(&admin).Error
	} else {
		var emp models.Employee
		if err = config.DB.First(&emp, userID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "员工不存在"})
			return
		}
		// 更新员工信息
		emp.Username = req.Name
		emp.Email = req.Email
		emp.Phone = req.Phone
		emp.Avatar = req.Avatar
		err = config.DB.Save(&emp).Error
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败: " + err.Error()})
		return
	}

	// 返回更新后的信息（复用 ToProfileResponse）
	c.JSON(http.StatusOK, gin.H{"message": "资料更新成功"})
}

func UpdatePassword(c *gin.Context) {
	// 从上下文中获取用户身份
	userID, _ := utils.GetCurrentUserID(c)
	role, _ := utils.GetCurrentUserRole(c)

	// 绑定请求参数
	var req models.UpdatePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误: " + err.Error()})
		return
	}

	// 根据角色查询用户
	var user models.UserPasswordHolder // 假设有一个公共接口或结构体包含密码字段
	var dbQuery *gorm.DB

	if role == "admin" {
		var admin models.Admin
		dbQuery = config.DB.Where("admin_id = ?", userID).First(&admin)
		user = &admin // 假设Admin实现了密码接口
	} else {
		var emp models.Employee
		dbQuery = config.DB.Where("emp_id = ?", userID).First(&emp)
		user = &emp // 假设Employee实现了密码接口
	}

	if dbQuery.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	// 验证旧密码是否正确
	if err := bcrypt.CompareHashAndPassword([]byte(user.GetPassword()), []byte(req.OldPassword)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "旧密码错误"})
		return
	}

	// 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密码加密失败"})
		return
	}

	// 更新数据库
	user.SetPassword(string(hashedPassword))
	if err := config.DB.Save(user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密码更新失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "密码修改成功"})
}

// controllers/auth.go
func Logout(c *gin.Context) {
	// 获取 Token
	authHeader := c.GetHeader("Authorization")
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	// 安全获取 claims（已由中间件确保存在）
	claims := c.MustGet("claims").(*utils.Claims)

	// 计算 Token 剩余有效期
	if claims.ExpiresAt == nil {
		c.JSON(401, gin.H{"error": "令牌缺少过期时间"})
		return
	}
	remaining := time.Until(claims.ExpiresAt.Time)

	// 将 Token 加入黑名单
	if err := config.Rdb.Set(config.Ctx, "jwt_blacklist:"+tokenString, "1", remaining).Err(); err != nil {
		c.JSON(500, gin.H{"error": "注销失败"})
		return
	}

	// 发送注销日志消息
	logData := map[string]interface{}{
		"user_id":   claims.UserID,
		"action":    "logout",
		"target_id": "", // 注销操作无目标对象
	}
	services.SendLogToRabbitMQ(logData)

	c.JSON(200, gin.H{"message": "已成功注销"})
}
