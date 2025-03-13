// Package models models/leave_request.go
package models

import "time"

type LeaveRequest struct {
	ID        uint       `gorm:"primaryKey;autoIncrement;column:id"`
	AdminID   *uint      `gorm:"column:admin_id;index"` // 使用指针，允许为空
	EmpID     *uint      `gorm:"column:emp_id;index"`   // 使用指针，允许为空
	Reason    string     `gorm:"type:text;size:100;not null"`
	StartTime *time.Time `gorm:"type:datetime;not null"`
	EndTime   *time.Time `gorm:"type:datetime;not null"`
	Status    string     `gorm:"type:enum('pending','approved','rejected');default:'pending'"`

	// 明确指定外键关系（重要修正）
	Admin    *Admin    `gorm:"foreignKey:AdminID;references:AdminID"`
	Employee *Employee `gorm:"foreignKey:EmpID;references:EmpID"`
}

// TableName models/leave_request.go
func (LeaveRequest) TableName() string {
	return "leave_requests" // 确保与数据库表名一致
}
