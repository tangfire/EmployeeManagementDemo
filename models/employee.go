package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Employee struct {
	EmpID     uint           `gorm:"primaryKey;autoIncrement;column:emp_id"`
	DepID     uint           `gorm:"column:dep_id"`
	Username  string         `gorm:"type:varchar(20);not null;unique;column:username"`
	Password  string         `gorm:"type:varchar(200);not null"`
	Position  string         `gorm:"type:varchar(50)"` // 职位
	Gender    string         `gorm:"type:varchar(10)"`
	Email     string         `gorm:"type:varchar(50)"`
	Phone     string         `gorm:"type:varchar(11)"`
	Avatar    string         `gorm:"type:varchar(100)"`
	Address   string         `gorm:"type:varchar(10)"`
	Salary    int            `gorm:"type:int(10)"`
	Status    string         `gorm:"type:varchar(20);default:'在职'"` // 状态（在职/离职）
	DeletedAt gorm.DeletedAt `gorm:"index"`                         // 软删除字段
}

func (Employee) TableName() string {
	return "employee"
}

// models/employee.go
func (e *Employee) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(e.Password), []byte(password))
	return err == nil
}

// models/employee.go
func (e *Employee) GetID() uint {
	return e.EmpID
}

func (e *Employee) GetUsername() string {
	return e.Username
}

func (e *Employee) GetRole() string {
	return "employee" // 标识角色为普通员工
}

// models/employee.go
type EmployeeProfileResponse struct {
	EmpID    uint   `json:"emp_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	DepID    uint   `json:"dep_id"`
}

func (e *Employee) ToProfileResponse() EmployeeProfileResponse {
	return EmployeeProfileResponse{
		EmpID:    e.EmpID,
		Username: e.Username,
		Email:    e.Email,
		Phone:    e.Phone,
		DepID:    e.DepID,
	}
}

// models/employee.go
func (e *Employee) GetPassword() string {
	return e.Password
}

func (e *Employee) SetPassword(pwd string) {
	e.Password = pwd
}
