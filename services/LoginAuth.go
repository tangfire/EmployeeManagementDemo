// services/LoginAuth.go
package services

import (
	"EmployeeManagementDemo/dao"
	"EmployeeManagementDemo/models"
	"errors"
)

// 认证逻辑
func AuthenticateUser(username, password string) (models.BaseUser, error) {
	// 先尝试查找管理员
	admin, err := dao.GetAdminByUsername(username)
	if err == nil && admin.CheckPassword(password) {
		return admin, nil
	}

	// 再尝试查找员工
	emp, err := dao.GetEmployeeByUsername(username)
	if err == nil && emp.CheckPassword(password) {
		return emp, nil
	}

	return nil, errors.New("invalid credentials")
}
