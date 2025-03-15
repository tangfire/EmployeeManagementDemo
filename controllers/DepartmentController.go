package controllers

import (
	"EmployeeManagementDemo/config"
	"EmployeeManagementDemo/models"
	"EmployeeManagementDemo/services"
	"EmployeeManagementDemo/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CreateDepartment(c *gin.Context) {
	var req models.CreateDepartmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误: " + err.Error()})
		return
	}

	// 检查部门名称是否已存在
	var existing models.Department
	if err := config.DB.Where("depart = ?", req.Depart).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "部门名称已存在"})
		return
	}

	// 创建新部门
	department := models.Department{Depart: req.Depart}
	if err := config.DB.Create(&department).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "部门创建成功",
		"dep_id":  department.DepID,
	})
}

func GetDepartments(c *gin.Context) {
	var departments []models.Department
	if err := config.DB.Find(&departments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取部门列表失败"})
		return
	}

	c.JSON(http.StatusOK, departments)
}

func UpdateDepartment(c *gin.Context) {
	depID := c.Param("dep_id") // 从URL路径获取部门ID

	var req models.UpdateDepartmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误: " + err.Error()})
		return
	}

	// 检查目标部门是否存在
	var department models.Department
	if err := config.DB.First(&department, depID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "部门不存在"})
		return
	}

	// 检查新名称是否冲突（排除自身）
	if req.Depart != department.Depart {
		var existing models.Department
		if err := config.DB.Where("depart = ? AND dep_id != ?", req.Depart, depID).First(&existing).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "部门名称已存在"})
			return
		}
	}

	// 更新字段
	department.Depart = req.Depart
	if err := config.DB.Save(&department).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "部门更新成功"})
}

// DeleteDepartment godoc
// @Summary 删除部门
// @Description 根据部门ID删除部门
// @Tags 部门管理
// @Accept json
// @Produce json
// @Param id path int true "部门ID"
// @Security BearerAuth
// @Success 200 {object} map[string]string "成功响应"
// @Failure 401 {object} map[string]string "未授权"
// @Failure 500 {object} map[string]string "内部错误"
// @Router /departments/{id} [delete]
func DeleteDepartment(c *gin.Context) {
	depID := c.Param("dep_id")
	adminId, err := utils.GetCurrentUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "请先登录!"})
		return
	}

	// 检查是否有员工关联
	var empCount int64
	config.DB.Model(&models.Employee{}).Where("department_id = ?", depID).Count(&empCount)
	if empCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "部门下存在员工，无法删除"})
		return
	}

	// 执行删除
	if err := config.DB.Delete(&models.Department{}, depID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败: " + err.Error()})
		return
	}

	// 发送操作日志
	services.SendLogToRabbitMQ(map[string]interface{}{
		"user_id":   adminId,
		"action":    "delete_department",
		"target_id": depID, // 或更详细的描述
	})

	c.JSON(http.StatusOK, gin.H{"message": "部门删除成功"})
}

// controllers/department_controller.go
func GetDepartmentSalaryAverages(c *gin.Context) {
	data, err := services.GetDepartmentAvgSalaries()
	if err != nil {
		c.JSON(500, models.Error(500, "获取数据失败"))
		return
	}

	c.JSON(200, models.Success(data))
}

// GetDepartmentHeadcounts controllers/department_controller.go
func GetDepartmentHeadcounts(c *gin.Context) {
	data, err := services.GetDepartmentHeadcounts()
	if err != nil {
		c.JSON(500, models.Error(500, "获取数据失败"))
		return
	}
	c.JSON(200, models.Success(data))
}
