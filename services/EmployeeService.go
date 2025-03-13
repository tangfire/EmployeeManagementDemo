package services

import (
	"EmployeeManagementDemo/dao"
	"EmployeeManagementDemo/models"
	"errors"
)

// services/LoginAuth.go
import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashed), err
}

// services/employee.go
func CreateEmployee(emp *models.Employee) error {
	// 检查用户名唯一性
	if exists, _ := dao.CheckUsernameExists(emp.Username); exists {
		return errors.New("用户名已存在")
	}

	// 默认部门分配（示例）
	//emp.DepID = 1 // 默认部门ID
	return dao.CreateEmployee(emp)
}
