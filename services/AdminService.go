package services

import (
	"EmployeeManagementDemo/dao"
	"EmployeeManagementDemo/models"
)

// ValidateAdminSecret 验证管理员密钥
func ValidateAdminSecret(secret string) bool {
	// 从配置读取密钥（示例值，实际应从安全配置加载）
	const validSecret = "tangfire"
	return secret == validSecret
}

// CheckAdminNameExists 检查管理员用户名是否存在
func CheckAdminNameExists(username string) (bool, error) {
	return dao.CheckAdminNameExists(username)
}

// CreateAdmin 创建管理员
func CreateAdmin(admin *models.Admin) error {
	return dao.CreateAdmin(admin)
}
