package controllers

import (
	"EmployeeManagementDemo/config"
	"EmployeeManagementDemo/models"
	"EmployeeManagementDemo/services"
	"EmployeeManagementDemo/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

// AdminRegister 管理员注册接口
func AdminRegister(c *gin.Context) {
	var req models.AdminRegisterRequest

	// 参数绑定验证
	if err := c.ShouldBindJSON(&req); err != nil {
		//c.JSON(http.StatusBadRequest, gin.H{"error": utils.TranslateValidationErrors(err)})
		c.JSON(http.StatusBadRequest, models.Error(400, utils.TranslateValidationErrors(err)))
		return
	}

	// 验证管理员密钥
	if valid := services.ValidateAdminSecret(req.SecretKey); !valid {
		//c.JSON(http.StatusForbidden, gin.H{"error": "无效的管理员密钥"})
		c.JSON(http.StatusForbidden, models.Error(400, "无效的管理员密钥"))
		return
	}

	// 检查用户名是否存在
	if exists, _ := services.CheckAdminNameExists(req.Username); exists {
		//c.JSON(http.StatusConflict, gin.H{"error": "管理员用户名已存在"})
		c.JSON(http.StatusConflict, models.Error(400, "管理员用户名已存在"))

		return
	}

	// 密码加密
	hashedPassword, err := services.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error(500, "密码处理失败"))
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
		c.JSON(http.StatusInternalServerError, models.Error(500, "创建管理员失败"))
		return
	}

	// 返回创建成功响应（隐藏敏感信息）
	//c.JSON(http.StatusCreated, gin.H{
	//	"message": "管理员账号创建成功",
	//	"data": gin.H{
	//		"admin_id": newAdmin.AdminID,
	//		"username": newAdmin.AdminName,
	//		"email":    newAdmin.AdminEmail,
	//	},
	//})

	c.JSON(http.StatusOK, models.Success(newAdmin))
}

func CreateEmployee(c *gin.Context) {
	var req models.CreateEmployeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误: " + err.Error()})
		return
	}

	// 检查邮箱是否已存在
	var existing models.Employee
	if err := config.DB.Where("email = ?", req.Email).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "邮箱已存在"})
		return
	}

	// 检查部门是否存在
	var department models.Department
	if err := config.DB.First(&department, req.DepartmentID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "部门不存在"})
		return
	}

	// 创建员工记录
	employee := models.Employee{
		Username: req.Name,
		DepID:    req.DepartmentID,
		Position: req.Position,
		Email:    req.Email,
		Phone:    req.Phone,
		Status:   "在职", // 默认状态
	}

	if err := config.DB.Create(&employee).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "员工创建成功",
		"emp_id":  employee.EmpID,
	})
}

func GetEmployees(c *gin.Context) {
	// 分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	offset := (page - 1) * pageSize

	query := config.DB.Model(&models.Employee{}).Where("deleted_at IS NULL")

	// 部门筛选（支持多选）
	if depIDs := c.QueryArray("dep_id"); len(depIDs) > 0 {
		query = query.Where("dep_id IN (?)", depIDs)
	}

	// 性别筛选（支持多选）
	if genders := c.QueryArray("gender"); len(genders) > 0 {
		query = query.Where("gender IN (?)", genders)
	}

	// 状态筛选（支持多选）
	if statuses := c.QueryArray("status"); len(statuses) > 0 {
		query = query.Where("status IN (?)", statuses)
	}

	// 全局搜索
	if search := c.Query("search"); search != "" {
		query = query.Where(
			"username LIKE ? OR position LIKE ? OR phone LIKE ?",
			"%"+search+"%", "%"+search+"%", "%"+search+"%",
		)
	}

	// 排序处理
	if sortField := c.Query("sortField"); sortField != "" {
		order := sortField
		if sortOrder := c.Query("sortOrder"); sortOrder == "descend" {
			order += " DESC"
		} else {
			order += " ASC"
		}
		query = query.Order(order)
	}

	// 分页查询
	var total int64
	query.Count(&total)

	var employees []models.Employee
	if err := query.Offset(offset).Limit(pageSize).Find(&employees).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.Error(500, "查询失败"))
		return
	}

	c.JSON(http.StatusOK, models.Success(gin.H{
		"data":  employees,
		"total": total,
	}))
}

func UpdateEmployee(c *gin.Context) {
	empID := c.Param("emp_id")

	var req models.UpdateEmployeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误: " + err.Error()})
		return
	}

	// 查找目标员工
	var employee models.Employee
	if err := config.DB.First(&employee, empID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "员工不存在"})
		return
	}

	// 邮箱唯一性校验（如果更新了邮箱）
	if req.Email != "" && req.Email != employee.Email {
		var existing models.Employee
		if err := config.DB.Where("email = ?", req.Email).First(&existing).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "邮箱已存在"})
			return
		}
		employee.Email = req.Email
	}

	// 更新其他字段
	if req.Name != "" {
		employee.Username = req.Name
	}
	if req.DepartmentID != 0 {
		// 检查新部门是否存在
		var department models.Department
		if err := config.DB.First(&department, req.DepartmentID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "部门不存在"})
			return
		}
		employee.DepID = req.DepartmentID
	}
	if req.Position != "" {
		employee.Position = req.Position
	}
	if req.Phone != "" {
		employee.Phone = req.Phone
	}
	if req.Status != "" {
		employee.Status = req.Status
	}

	if err := config.DB.Save(&employee).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "员工信息更新成功"})
}

// DeleteEmployee godoc
// @Summary 删除员工
// @Description 根据员工ID删除员工
// @Tags 员工管理
// @Accept json
// @Produce json
// @Param id path int true "员工ID"
// @Security BearerAuth
// @Success 200 {object} map[string]string "成功响应"
// @Failure 401 {object} map[string]string "未授权"
// @Failure 500 {object} map[string]string "内部错误"
// @Router /employees/{id} [delete]
func DeleteEmployee(c *gin.Context) {
	empID := c.Param("emp_id")
	adminId, err := utils.GetCurrentUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "请先登录!"})
		return
	}

	// 执行删除（硬删除，如需软删除需修改模型）
	if err := config.DB.Delete(&models.Employee{}, empID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败: " + err.Error()})
		return
	}

	// 发送操作日志
	services.SendLogToRabbitMQ(map[string]interface{}{
		"user_id":   adminId,
		"action":    "delete_employee",
		"target_id": empID,
	})

	c.JSON(http.StatusOK, gin.H{"message": "员工删除成功"})
}

func ApproveLeaveRequest(c *gin.Context) {
	// 获取管理员ID
	adminID, err := utils.GetCurrentUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "请先登录"})
		return
	}

	leaveID := c.Param("id")
	var req models.ApproveLeaveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误: " + err.Error()})
		return
	}

	// 查找待审批的请假记录
	var leave models.LeaveRequest
	if err := config.DB.First(&leave, "id = ? AND status = 'pending'", leaveID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "未找到待审批的申请"})
		return
	}

	// 更新审批状态和管理员ID
	leave.AdminID = &adminID
	leave.Status = req.Status
	if err := config.DB.Save(&leave).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "审批失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "审批结果已提交"})
}

func GetAllLeaveRequests(c *gin.Context) {
	status := c.Query("status") // 支持按状态过滤
	query := config.DB.Model(&models.LeaveRequest{})

	if status != "" {
		query = query.Where("status = ?", status)
	}

	var leaves []models.LeaveRequest
	query.Preload("Employee").Preload("Admin").Find(&leaves) // 预加载关联数据

	c.JSON(http.StatusOK, leaves)
}

// controllers/admin.go
func KickUser(c *gin.Context) {
	userID := c.Param("user_id")
	adminID, err2 := utils.GetCurrentUserID(c)
	if err2 != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "请先登录"})
		return
	}

	// 记录踢出时间戳
	invalidateKey := "user_invalid:" + userID
	err := config.Rdb.Set(config.Ctx, invalidateKey, time.Now().Unix(), 0).Err()
	if err != nil {
		c.JSON(500, gin.H{"error": "操作失败"})
		return
	}

	//发送踢人日志消息
	logData := map[string]interface{}{
		"user_id":   adminID,
		"action":    "kick_user",
		"target_id": userID,
	}
	services.SendLogToRabbitMQ(logData)

	c.JSON(200, gin.H{"message": "用户已被踢出"})
}
