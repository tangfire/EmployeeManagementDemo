package models

import (
	"time"
)

type SignRecord struct {
	ID          uint       `gorm:"primaryKey;autoIncrement;column:id"`
	EmpID       uint       `gorm:"column:emp_id;index;not null"` // 使用指针
	SignInTime  *time.Time `gorm:"type:datetime;"`
	SignOutTime *time.Time `gorm:"type:datetime;"`
	Date        time.Time  `gorm:"index;type:date"` // 记录日期（用于快速查询当天记录）

	// 添加级联约束（根据业务需求）
	Employee *Employee `gorm:"foreignKey:EmpID;references:EmpID"`
}

func (SignRecord) TableName() string {
	return "sign_records"
}
