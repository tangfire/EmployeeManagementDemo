package controllers

import (
	"EmployeeManagementDemo/config"
	"EmployeeManagementDemo/models"
	"EmployeeManagementDemo/services"
	"EmployeeManagementDemo/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// controllers/LoginAuth.go
func Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.TranslateValidationErrors(err))
		return
	}

	// 密码加密
	hashedPassword, err := services.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密码加密失败"})
		return
	}

	// 创建员工账号
	emp := models.Employee{
		Username: req.Username,
		Password: hashedPassword,
		Email:    req.Email,
		Phone:    req.Phone,
	}

	if err := services.CreateEmployee(&emp); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "用户名已存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "注册成功"})
}

func SignIn(c *gin.Context) {
	// 获取当前用户ID（假设只有员工可以签到）
	userID, err := utils.GetCurrentUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未登录"})
		return
	}

	// 检查当天是否已签到
	var existingRecord models.SignRecord
	today := time.Now().Format("2006-01-02")
	result := config.DB.Where("emp_id = ? AND date = ?", userID, today).First(&existingRecord)

	if result.Error == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "今日已签到"})
		return
	}

	// 创建签到记录
	now := time.Now()
	newRecord := models.SignRecord{
		EmpID:      userID,
		SignInTime: &now,
		Date:       now,
	}

	if err := config.DB.Create(&newRecord).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "签到失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "签到成功",
		"sign_in_time": now.Format(time.RFC3339),
	})
}

func SignOut(c *gin.Context) {
	userID, err := utils.GetCurrentUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未登录"})
		return
	}

	// 查找当天未签退的记录
	var record models.SignRecord
	today := time.Now().Format("2006-01-02")
	result := config.DB.Where("emp_id = ? AND date = ? AND sign_out_time IS NULL", userID, today).First(&record)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "未找到有效签到记录"})
		return
	}

	// 更新签退时间
	now := time.Now()
	record.SignOutTime = &now
	if err := config.DB.Save(&record).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "签退失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "签退成功",
		"sign_out_time": now.Format(time.RFC3339),
	})
}

func CreateLeaveRequest(c *gin.Context) {
	// 获取当前登录员工ID
	empID, err := utils.GetCurrentUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "请先登录"})
		return
	}

	var req models.CreateLeaveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误: " + err.Error()})
		return
	}

	// 创建请假记录
	leave := models.LeaveRequest{
		EmpID:     &empID,
		Reason:    req.Reason,
		StartTime: &req.StartTime,
		EndTime:   &req.EndTime,
		Status:    "pending",
	}

	if err := config.DB.Create(&leave).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "提交失败"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "请假申请已提交",
		"id":      leave.ID,
	})
}

func GetMyLeaveRequests(c *gin.Context) {
	empID, err := utils.GetCurrentUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "请先登录"})
		return
	}

	var leaves []models.LeaveRequest
	config.DB.Where("emp_id = ?", empID).Find(&leaves)

	c.JSON(http.StatusOK, leaves)
}
