package controllers

import (
	"EmployeeManagementDemo/config"
	"EmployeeManagementDemo/models"
	"EmployeeManagementDemo/services"
	"EmployeeManagementDemo/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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

	//query := config.DB.Model(&models.Employee{}).Where("deleted_at IS NULL")
	// 修改后
	query := config.DB.
		Model(&models.Employee{}).
		Select("employees.*, departments.depart as dep_name").
		Joins("LEFT JOIN departments ON employees.dep_id = departments.dep_id").
		Where("employees.deleted_at IS NULL")

	// 部门筛选（支持多选）
	if depIDs := c.QueryArray("dep_id"); len(depIDs) > 0 {
		query = query.Where("dep_id IN (?)", depIDs)
	}

	// 性别筛选（支持多选）
	// 兼容两种参数格式（gender[]=男）

	if genders := c.QueryArray("gender[]"); len(genders) > 0 {
		query = query.Where("gender IN (?)", genders)
	}

	// 状态筛选（支持多选）
	if statuses := c.QueryArray("status[]"); len(statuses) > 0 {
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

	var employeesWithDepNameDto []models.EmployeeWithDepNameDTO
	if err := query.Offset(offset).Limit(pageSize).Find(&employeesWithDepNameDto).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.Error(500, "查询失败"))
		return
	}

	c.JSON(http.StatusOK, models.Success(gin.H{
		"data":  employeesWithDepNameDto,
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

// ExportEmployees 导出接口(事务版)
func ExportEmployees(c *gin.Context) {
	// 开启事务（隔离级别设为REPEATABLE READ）
	tx := config.DB.Begin()
	if tx.Error != nil {
		c.JSON(500, models.Error(500, "事务启动失败"))
		return
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 事务内查询（网页1][3]
	var employeesWithDepNameDto []models.EmployeeWithDepNameDTO
	query := tx.Model(&models.Employee{}).
		Select("employees.*, departments.depart as dep_name").
		Joins("LEFT JOIN departments ON employees.dep_id = departments.dep_id").
		Where("employees.deleted_at IS NULL")

	if err := query.Find(&employeesWithDepNameDto).Error; err != nil {
		tx.Rollback()
		c.JSON(500, models.Error(500, "数据查询失败"))
		return
	}

	// 创建Excel文件
	f := excelize.NewFile()
	sheet := "员工信息"
	index, _ := f.NewSheet(sheet)
	f.SetActiveSheet(index)

	// 设置表头
	headers := []string{"工号", "姓名", "部门", "职位", "性别", "薪资", "状态"}
	for col, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(col+1, 1)
		f.SetCellValue(sheet, cell, h)
	}

	// 填充数据
	for row, emp := range employeesWithDepNameDto {
		rowIndex := row + 2
		f.SetCellValue(sheet, fmt.Sprintf("A%d", rowIndex), emp.EmpID)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", rowIndex), emp.Username)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", rowIndex), emp.DepName)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", rowIndex), emp.Position)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", rowIndex), emp.Gender)
		f.SetCellValue(sheet, fmt.Sprintf("F%d", rowIndex), emp.Salary)
		f.SetCellValue(sheet, fmt.Sprintf("G%d", rowIndex), emp.Status)
	}

	// 提交事务（网页2][3]
	if err := tx.Commit().Error; err != nil {
		c.JSON(500, models.Error(500, "事务提交失败"))
		return
	}

	// 设置响应头
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=employees.xlsx")

	// 输出文件流
	if _, err := f.WriteTo(c.Writer); err != nil {
		c.JSON(500, models.Error(500, "文件流输出失败"))
	}
}

// ExportEmployees 导出接口(无事务)
//func ExportEmployees(c *gin.Context) {
//	// 获取所有数据
//	var employeesWithDepNameDto []models.EmployeeWithDepNameDTO
//	//if err := config.DB.Find(&employees).Error; err != nil {
//	//	c.JSON(500, models.Error(500, "数据查询失败"))
//	//	return
//	//}
//	query := config.DB.
//		Model(&models.Employee{}).
//		Select("employees.*, departments.depart as dep_name").
//		Joins("LEFT JOIN departments ON employees.dep_id = departments.dep_id").
//		Where("employees.deleted_at IS NULL")
//
//	if err := query.Find(&employeesWithDepNameDto).Error; err != nil {
//		c.JSON(http.StatusInternalServerError, models.Error(500, "查询失败"))
//		return
//	}
//
//	// 创建Excel文件
//	f := excelize.NewFile()
//	sheet := "员工信息"
//	index, _ := f.NewSheet(sheet)
//	f.SetActiveSheet(index) // 添加此行
//
//	// 设置表头
//	headers := []string{"工号", "姓名", "部门", "职位", "性别", "薪资", "状态"}
//	for col, h := range headers {
//		cell, _ := excelize.CoordinatesToCellName(col+1, 1)
//		f.SetCellValue(sheet, cell, h)
//	}
//
//	// 填充数据
//	for row, emp := range employeesWithDepNameDto {
//		rowIndex := row + 2
//		f.SetCellValue(sheet, fmt.Sprintf("A%d", rowIndex), emp.EmpID)
//		f.SetCellValue(sheet, fmt.Sprintf("B%d", rowIndex), emp.Username)
//		f.SetCellValue(sheet, fmt.Sprintf("C%d", rowIndex), emp.DepName)
//		f.SetCellValue(sheet, fmt.Sprintf("D%d", rowIndex), emp.Position)
//		f.SetCellValue(sheet, fmt.Sprintf("E%d", rowIndex), emp.Gender)
//		f.SetCellValue(sheet, fmt.Sprintf("F%d", rowIndex), emp.Salary)
//		f.SetCellValue(sheet, fmt.Sprintf("G%d", rowIndex), emp.Status)
//	}
//
//	// 设置响应头
//	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
//	c.Header("Content-Disposition", "attachment; filename=employees.xlsx")
//
//	// 输出文件流
//	if _, err := f.WriteTo(c.Writer); err != nil {
//		c.JSON(500, models.Error(500, "文件生成失败"))
//	}
//}

// ImportEmployees 导入接口(自动事务模式)
func ImportEmployees(c *gin.Context) {
	// 文件处理部分保持不变
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(400, models.Error(400, "文件上传失败"))
		return
	}

	// 校验文件格式
	if !strings.HasSuffix(file.Filename, ".xlsx") {
		c.JSON(400, models.Error(400, "仅支持.xlsx格式"))
		return
	}

	// 创建上传目录
	if err := os.MkdirAll("./uploads", os.ModePerm); err != nil {
		c.JSON(500, models.Error(500, "服务器存储目录创建失败"))
		return
	}

	// 保存文件
	dstPath := filepath.Join("./uploads", file.Filename)
	if err := c.SaveUploadedFile(file, dstPath); err != nil {
		c.JSON(400, models.Error(400, "文件保存失败: "+err.Error()))
		return
	}

	// 核心事务逻辑
	err = config.DB.Transaction(func(tx *gorm.DB) error {
		f, err := excelize.OpenFile(dstPath)
		if err != nil {
			return fmt.Errorf("文件格式错误: %v", err)
		}

		rows, _ := f.GetRows("员工信息")
		for i, row := range rows {
			if i == 0 {
				continue // 跳过表头
			}

			// 工号转换
			empID, err := strconv.ParseUint(row[0], 10, 64)
			if err != nil {
				return fmt.Errorf("第%d行工号格式错误", i+1)
			}

			// 薪资转换
			salaryValue, err := strconv.ParseFloat(row[5], 64)
			if err != nil {
				return fmt.Errorf("第%d行薪资格式错误", i+1)
			}

			// 部门查询（使用事务对象）
			depID, err := getDepIDByNameWithTx(tx, row[2])
			if err != nil {
				return fmt.Errorf("第%d行部门不存在", i+1)
			}

			emp := models.Employee{
				EmpID:    uint(empID),
				Username: row[1],
				DepID:    depID,
				Position: row[3],
				Gender:   row[4],
				Salary:   salaryValue,
				Status:   row[6],
			}

			// 数据校验（使用事务对象）
			if err := validateEmployeeWithTx(tx, emp); err != nil {
				return fmt.Errorf("第%d行数据错误: %v", i+1, err)
			}

			// 数据库写入（使用事务对象）
			if err := tx.Create(&emp).Error; err != nil {
				return fmt.Errorf("第%d行保存失败: %v", i+1, err)
			}
		}
		return nil // 全部成功自动提交
	})

	// 统一错误处理
	if err != nil {
		c.JSON(400, models.Error(400, err.Error()))
		return
	}

	c.JSON(200, models.Success(nil))
}

// 改造后的部门查询函数(支持事务)
func getDepIDByNameWithTx(tx *gorm.DB, name string) (uint, error) {
	var dep models.Department
	if err := tx.Where("depart = ?", name).First(&dep).Error; err != nil {
		return 0, err
	}
	return dep.DepID, nil
}

// 改造后的数据校验函数(支持事务)
func validateEmployeeWithTx(tx *gorm.DB, emp models.Employee) error {
	// 基础校验
	if emp.EmpID == 0 {
		return fmt.Errorf("工号不能为空")
	}
	if emp.Username == "" {
		return fmt.Errorf("姓名不能为空")
	}

	// 部门存在性校验（使用事务对象）
	var dep models.Department
	if err := tx.Where("dep_id = ?", emp.DepID).First(&dep).Error; err != nil {
		return fmt.Errorf("部门ID不存在")
	}

	// 枚举值校验
	validGenders := map[string]bool{"男": true, "女": true, "其他": true}
	if !validGenders[emp.Gender] {
		return fmt.Errorf("性别值无效")
	}

	// 唯一性校验（使用事务对象）
	var existing models.Employee
	if err := tx.Where("emp_id = ? OR username = ?", emp.EmpID, emp.Username).
		First(&existing).Error; err == nil {
		return fmt.Errorf("工号或用户名已存在")
	}

	return nil
}

//// ImportEmployees 导入接口(无事务)
//func ImportEmployees(c *gin.Context) {
//	file, err := c.FormFile("file")
//	if err != nil {
//		c.JSON(400, models.Error(400, "文件上传失败"))
//		return
//	}
//
//	// 校验文件格式
//	if !strings.HasSuffix(file.Filename, ".xlsx") {
//		c.JSON(400, models.Error(400, "仅支持.xlsx格式"))
//		return
//	}
//
//	// 创建上传目录
//	if err := os.MkdirAll("./uploads", os.ModePerm); err != nil {
//		c.JSON(500, models.Error(500, "服务器存储目录创建失败"))
//		return
//	}
//
//	// 保存文件
//	dstPath := filepath.Join("./uploads", file.Filename) // 使用跨平台路径拼接
//	if err := c.SaveUploadedFile(file, dstPath); err != nil {
//		c.JSON(400, models.Error(400, "文件保存失败: "+err.Error()))
//		return
//	}
//
//	// 打开Excel文件（使用保存后的路径）
//	f, err := excelize.OpenFile(dstPath)
//	if err != nil {
//		c.JSON(400, models.Error(400, "文件格式错误: "+err.Error()))
//		return
//	}
//
//	// 解析数据
//	rows, _ := f.GetRows("员工信息")
//	for i, row := range rows {
//		if i == 0 {
//			continue
//		}
//
//		// 工号转换
//		empID, err := strconv.ParseUint(row[0], 10, 64)
//		if err != nil {
//			c.JSON(400, models.Error(400, fmt.Sprintf("第%d行工号格式错误", i+1)))
//			return
//		}
//
//		// 薪资转换
//		salaryValue, err := strconv.ParseFloat(row[5], 64)
//		if err != nil {
//			c.JSON(400, models.Error(400, fmt.Sprintf("第%d行薪资格式错误", i+1)))
//			return
//		}
//
//		// 部门ID转换
//		depID, err := getDepIDByName(row[2])
//		if err != nil {
//			c.JSON(400, models.Error(400, fmt.Sprintf("第%d行部门不存在", i+1)))
//			return
//		}
//
//		emp := models.Employee{
//			EmpID:    uint(empID),
//			Username: row[1],
//			DepID:    depID,
//			Position: row[3],
//			Gender:   row[4],
//			Salary:   salaryValue,
//			Status:   row[6],
//		}
//
//		// 数据校验
//		if err := validateEmployee(emp); err != nil {
//			c.JSON(400, models.Error(400, fmt.Sprintf("第%d行数据错误: %v", i+1, err)))
//			return
//		}
//
//		// 保存到数据库
//		if err := config.DB.Create(&emp).Error; err != nil {
//			c.JSON(500, models.Error(500, "数据保存失败"))
//			return
//		}
//	}
//
//	c.JSON(200, models.Success(nil))
//}
//
//// 优化后的部门查询函数(无事务)
//func getDepIDByName(name string) (uint, error) {
//	var dep models.Department
//	result := config.DB.Where("depart = ?", name).First(&dep)
//	if result.Error != nil {
//		return 0, result.Error
//	}
//	return dep.DepID, nil
//}
//
//// models/employee.go(无事务)
//func validateEmployee(emp models.Employee) error {
//	// 基础校验
//	if emp.EmpID == 0 {
//		return fmt.Errorf("工号不能为空")
//	}
//	if emp.Username == "" {
//		return fmt.Errorf("姓名不能为空")
//	}
//
//	// 部门存在性校验
//	var dep models.Department
//	if err := config.DB.Where("dep_id = ?", emp.DepID).First(&dep).Error; err != nil {
//		return fmt.Errorf("部门ID不存在")
//	}
//
//	// 枚举值校验
//	validGenders := map[string]bool{"男": true, "女": true, "其他": true}
//	if !validGenders[emp.Gender] {
//		return fmt.Errorf("性别值无效")
//	}
//
//	return nil
//}
