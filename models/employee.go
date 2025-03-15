package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Employee struct {
	EmpID     uint           `gorm:"primaryKey;autoIncrement;column:emp_id"`
	DepID     uint           `gorm:"column:dep_id;index;comment:所属部门ID"`
	Username  string         `gorm:"type:varchar(20);not null;unique;column:username"`
	Password  string         `gorm:"type:varchar(200);not null"`
	Position  string         `gorm:"type:varchar(50)"` // 职位
	Gender    string         `gorm:"type:enum('男','女','其他');default:'其他'"`
	Email     string         `gorm:"type:varchar(50)"`
	Phone     string         `gorm:"type:char(11);not null"` // 手机号必填
	Avatar    string         `gorm:"type:varchar(100)"`
	Address   string         `gorm:"type:varchar(100)"` // 地址长度扩展至100
	Salary    int            `gorm:"type:int(10)"`
	Status    string         `gorm:"type:varchar(20);default:'在职';index:idx_emp_status"` // 状态（在职/离职）
	DeletedAt gorm.DeletedAt `gorm:"index"`                                              // 软删除字段
}

func (Employee) TableName() string {
	return "employees"
}

// CheckPassword models/employee.go
func (e *Employee) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(e.Password), []byte(password))
	return err == nil
}

// GetID models/employee.go
func (e *Employee) GetID() uint {
	return e.EmpID
}

func (e *Employee) GetUsername() string {
	return e.Username
}

func (e *Employee) GetRole() string {
	return "employee" // 标识角色为普通员工
}

// EmployeeProfileResponse models/employee.go
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
