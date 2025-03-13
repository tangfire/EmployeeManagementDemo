// Package models models/sign_record.go
package models

import (
	"time"
)

type SignRecord struct {
	ID          uint      `gorm:"primaryKey;autoIncrement;column:id"` // 主键自增
	EmpID       uint      `gorm:"column:emp_id"`                      // 外键（关联员工表）
	SignInTime  time.Time `gorm:"type:datetime;not null"`             // 签到时间（非空）
	SignOutTime time.Time `gorm:"type:datetime;not null"`             // 签退时间（非空）

	// 关联模型（按需添加）
	Employee Employee `gorm:"foreignKey:EmpID;references:EmpID"` // 关联员工表
}

// TableName 自定义表名（与数据库表名一致）
func (SignRecord) TableName() string {
	return "sign_records"
}
