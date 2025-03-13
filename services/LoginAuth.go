// services/LoginAuth.go
package services

import (
	"EmployeeManagementDemo/dao"
	"errors"
)

// 定义用户接口（统一Admin和Employee的认证行为）
type User interface {
	GetID() uint
	GetUsername() string
	GetRole() string
	CheckPassword(password string) bool
}

// 认证逻辑
func AuthenticateUser(username, password string) (User, error) {
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
