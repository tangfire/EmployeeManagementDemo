// Package models models/employee.go
package models

type Employee struct {
	EmpID    uint   `gorm:"primaryKey;autoIncrement;column:emp_id"` // 主键自增
	DepID    uint   `gorm:"column:dep_id"`                          // 外键（需关联部门表）
	Username string `gorm:"type:varchar(20);not null;unique"`       // 不能为空且唯一
	Password string `gorm:"type:varchar(20);not null"`              // 不能为空
	Gender   string `gorm:"type:varchar(10)"`
	Email    string `gorm:"type:varchar(50)"`
	Phone    string `gorm:"type:varchar(11)"`
	Avatar   string `gorm:"type:varchar(100)"`
	Address  string `gorm:"type:varchar(10)"`
	Salary   int    `gorm:"type:int(10)"`
}

// TableName 自定义表名（可选）
func (Employee) TableName() string {
	return "employee"
}
