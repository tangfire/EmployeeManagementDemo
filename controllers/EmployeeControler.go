package controllers

import (
	"EmployeeManagementDemo/config"
	"EmployeeManagementDemo/models"
	"EmployeeManagementDemo/services"
	"EmployeeManagementDemo/utils"
	"database/sql"
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
	// 获取当前用户ID（示例值，需替换实际获取逻辑）
	userID, err := utils.GetCurrentUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.Error(400, "未登录"))
	}

	now := time.Now()
	today := now.Format("2006-01-02")

	// 检查是否重复签到
	var existingRecord models.SignRecord
	if result := config.DB.Where("emp_id = ? AND date = ?", userID, today).First(&existingRecord); result.Error == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "今日已签到"})
		return
	}

	// 判断是否迟到（示例规则：9:30前不算迟到）
	expectedSignIn := time.Date(now.Year(), now.Month(), now.Day(), 9, 30, 0, 0, now.Location())
	status := "正常"
	if now.After(expectedSignIn) {
		status = "迟到"
	}

	// 创建记录时修正日期精度
	newRecord := models.SignRecord{
		EmpID:      userID,
		SignInTime: now,
		Date:       now.Truncate(24 * time.Hour), // 仅保留日期
		Status:     status,
	}

	if err := config.DB.Create(&newRecord).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.Error(500, "签到失败"))
		return
	}

	c.JSON(http.StatusOK, models.Success(gin.H{
		"message":      "签到成功",
		"sign_in_time": now.Format(time.RFC3339),
		"status":       status,
	}))
}

func SignOut(c *gin.Context) {
	userID, err := utils.GetCurrentUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.Error(400, "未登录"))
	}
	today := time.Now().Format("2006-01-02")

	// 查找当天有效记录
	var record models.SignRecord
	result := config.DB.Where("emp_id = ? AND date = ? AND sign_out_time IS NULL", userID, today).First(&record)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, models.Error(400, "未找到有效签到记录"))
		return
	}

	// 更新签退时间和状态
	now := time.Now()
	// 修正签退时间赋值
	record.SignOutTime = sql.NullTime{Time: now, Valid: true}

	// 增强状态判断逻辑
	expectedSignOut := time.Date(now.Year(), now.Month(), now.Day(), 17, 0, 0, 0, now.Location())
	if now.Before(expectedSignOut) {
		if record.Status == "迟到" {
			record.Status = "迟到+早退"
		} else {
			record.Status = "早退"
		}
	} else if record.Status == "缺勤" { // 处理默认缺勤状态
		record.Status = "正常"
	}

	if err := config.DB.Save(&record).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.Error(500, "签退失败"))
		return
	}

	c.JSON(http.StatusOK, models.Success(gin.H{
		"message":       "签退成功",
		"sign_out_time": now.Format(time.RFC3339),
		"status":        record.Status,
	}))
}

func GetAttendance(c *gin.Context) {
	userID, err := utils.GetCurrentUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.Error(401, "用户未登录"))
		return
	}

	yearMonth := c.Query("month")
	if yearMonth == "" {
		c.JSON(http.StatusBadRequest, models.Error(400, "月份参数必填"))
		return
	}

	// 解析月份参数
	startTime, err := time.ParseInLocation("2006-01", yearMonth, time.Local)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error(400, "月份格式无效（正确示例：2025-03）"))
		return
	}

	// 设定日期范围
	loc := startTime.Location()
	startOfMonth := time.Date(startTime.Year(), startTime.Month(), 1, 0, 0, 0, 0, loc)
	endOfMonth := startOfMonth.AddDate(0, 1, 0).Add(-time.Nanosecond)

	// 格式化为字符串
	startStr := startOfMonth.Format("2006-01-02")
	endStr := endOfMonth.Format("2006-01-02")

	// 查询数据库
	var records []models.SignRecord
	err = config.DB.Debug().
		Where("emp_id = ? AND date BETWEEN ? AND ?", userID, startStr, endStr).
		Order("date DESC").
		Find(&records).
		Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error(500, "查询失败"))
		return
	}

	// 构造响应数据
	responseData := make([]map[string]interface{}, len(records))
	for i, record := range records {
		signOutTime := ""
		if record.SignOutTime.Valid {
			signOutTime = record.SignOutTime.Time.Format("15:04:05")
		}
		responseData[i] = map[string]interface{}{
			"date":          record.Date.Format("2006-01-02"),
			"sign_in_time":  record.SignInTime.Format("15:04:05"),
			"sign_out_time": signOutTime,
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
