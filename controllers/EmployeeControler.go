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
		c.JSON(http.StatusBadRequest, models.Error(400, utils.TranslateValidationErrors(err)))
		return
	}

	if req.Password != req.ConfirmPassword {
		c.JSON(http.StatusBadRequest, models.Error(400, "两次密码输入不相同"))
		return
	}

	// 密码加密
	hashedPassword, err := services.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error(500, "密码加密失败"))
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
		c.JSON(http.StatusConflict, models.Error(409, "用户名已存在"))
		return
	}

	//c.JSON(http.StatusOK, gin.H{"message": "注册成功"})
	// 注册成功时
	c.JSON(http.StatusOK, models.ApiResponse{
		Code:      200,
		Message:   "注册成功",
		Data:      emp,
		Timestamp: time.Now().Unix(),
	})
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
		SignInTime: now,
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
	record.SignOutTime = now
	if err := config.DB.Save(&record).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "签退失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "签退成功",
		"sign_out_time": now.Format(time.RFC3339),
	})
}

// GetAttendance 新增考勤查询接口（添加到原有签到签退路由旁）
func GetAttendance(c *gin.Context) {
	// 鉴权获取当前用户
	userID, err := utils.GetCurrentUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未登录"})
		return
	}

	// 处理查询参数
	yearMonth := c.Query("month") // 格式示例：2025-03
	if yearMonth == "" {
		c.JSON(http.StatusBadRequest, models.Error(400, "月份参数必填"))
		return
	}

	// 计算日期范围
	startTime, _ := time.Parse("2006-01", yearMonth)
	endTime := startTime.AddDate(0, 1, -1)

	// 查询数据库
	var records []models.SignRecord
	if err := config.DB.Where("emp_id = ? AND date BETWEEN ? AND ?",
		userID,
		startTime.Format("2006-01-02"),
		endTime.Format("2006-01-02"),
	).Order("date DESC").Find(&records).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.Error(500, "查询失败"))
		return
	}

	// 转换为响应格式
	responseData := make([]map[string]interface{}, len(records))
	for i, record := range records {
		responseData[i] = map[string]interface{}{
			"date":          record.Date.Format("2006-01-02"),
			"sign_in_time":  record.SignInTime.Format("15:04:05"),
			"sign_out_time": record.SignOutTime.Format("15:04:05"),
			"status":        record.Status,
		}
	}

	c.JSON(http.StatusOK, models.Success(responseData))
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
		EmpID:     empID,
		Reason:    req.Reason,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
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
