package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Employee struct {
	EmpID     uint           `gorm:"primaryKey;autoIncrement;column:emp_id" json:"emp_id"`
	DepID     uint           `gorm:"column:dep_id;index;comment:所属部门ID" json:"dep_id"`
	Username  string         `gorm:"type:varchar(20);not null;unique;column:username" json:"username"`
	Password  string         `gorm:"type:varchar(200);not null" json:"password"`
	Position  string         `gorm:"type:varchar(50)" json:"position"`
	Gender    string         `gorm:"type:enum('男','女','其他');default:'其他'" json:"gender"`
	Email     string         `gorm:"type:varchar(50)" json:"email"`
	Phone     string         `gorm:"type:char(11);not null" json:"phone"`
	Avatar    string         `gorm:"type:varchar(100)" json:"avatar"`
	Address   string         `gorm:"type:varchar(100)" json:"address"`
	Salary    float64        `gorm:"type:int(10)" json:"salary"`
	Status    string         `gorm:"type:varchar(20);default:'在职';index:idx_emp_status" json:"status"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
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
