package models

import "golang.org/x/crypto/bcrypt"

type Employee struct {
	EmpID    uint   `gorm:"primaryKey;autoIncrement;column:emp_id"`
	DepID    uint   `gorm:"column:dep_id"`
	Username string `gorm:"type:varchar(20);not null;unique;column:username"`
	Password string `gorm:"type:varchar(200);not null"`
	Gender   string `gorm:"type:varchar(10)"`
	Email    string `gorm:"type:varchar(50)"`
	Phone    string `gorm:"type:varchar(11)"`
	Avatar   string `gorm:"type:varchar(100)"`
	Address  string `gorm:"type:varchar(10)"`
	Salary   int    `gorm:"type:int(10)"`
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
