package dao

import (
	"EmployeeManagementDemo/config"
	"EmployeeManagementDemo/models"
	"log"
)

// dao/employee.go
func CheckUsernameExists(username string) (bool, error) {
	var count int64

	// 使用 Debug 模式直接打印 SQL
	result := config.DB.Debug().Model(&models.Employee{}).
		Where("username = ?", username).
		Count(&count)

	// 强制触发 SQL 构建（调试用）
	log.Printf("完整 SQL: %v", result.Statement.SQL.String())

	return count > 0, result.Error
}

func CreateEmployee(emp *models.Employee) error {
	return config.DB.Create(emp).Error
}

func GetEmployeeByUsername(username string) (*models.Employee, error) {
	var emp models.Employee
	result := config.DB.Where("username = ?", username).First(&emp)
	return &emp, result.Error
}
