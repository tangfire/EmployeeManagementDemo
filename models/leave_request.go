// Package models models/leave_request.go
package models

import (
	"time"
)

type LeaveRequest struct {
	ID        uint      `gorm:"primaryKey;autoIncrement;column:id"` // 主键自增
	AdminID   uint      `gorm:"column:admin_id"`                    // 外键（关联管理员表）
	EmpID     uint      `gorm:"column:emp_id"`                      // 外键（关联员工表）
	Reason    string    `gorm:"type:text;size:100;not null"`        // 非空，长度限制
	StartTime time.Time `gorm:"type:datetime;not null"`             // 非空时间
	EndTime   time.Time `gorm:"type:datetime;not null"`             // 非空时间
	Status    string    `gorm:"type:enum('pending','approved','rejected');default:'pending'"`

	// 关联模型（按需定义）
	Admin Admin    `gorm:"foreignKey:AdminID"` // 关联管理员
	User  Employee `gorm:"foreignKey:EmpID"`   // 关联员工
}

// TableName 自定义表名
func (LeaveRequest) TableName() string {
	return "leave_requests"
}
