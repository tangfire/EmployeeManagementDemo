package dao

import (
	"EmployeeManagementDemo/config"
	"EmployeeManagementDemo/models"
)

func CheckAdminNameExists(username string) (bool, error) {
	var count int64
	result := config.DB.Model(&models.Admin{}).
		Where("admin_name = ?", username).
		Count(&count)
	return count > 0, result.Error
}

func CreateAdmin(admin *models.Admin) error {
	return config.DB.Create(admin).Error
}

func GetAdminByUsername(username string) (*models.Admin, error) {
	var admin models.Admin
	result := config.DB.Where("admin_name = ?", username).First(&admin)
	return &admin, result.Error
}
