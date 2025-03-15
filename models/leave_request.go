// Package models models/leave_request.go
package models

import (
	"gorm.io/gorm"
	"math"
	"time"
)

type LeaveRequest struct {
	ID         uint      `gorm:"primaryKey;autoIncrement;column:id"`
	AdminID    *uint     `gorm:"column:admin_id;index"`        // 使用指针，允许为空
	EmpID      uint      `gorm:"column:emp_id;index;not null"` // 使用指针，允许为空
	Reason     string    `gorm:"type:text;size:100;not null"`
	StartTime  time.Time `gorm:"type:datetime;not null"` // 非指针
	EndTime    time.Time `gorm:"type:datetime;not null"`
	ApprovedAt *time.Time
	Status     string  `gorm:"type:enum('pending','approved','rejected');default:'pending';index:idx_status"`
	Duration   float64 `gorm:"column:duration;comment:请假天数"` // 新增字段

}

// TableName models/leave_request.go
func (LeaveRequest) TableName() string {
	return "leave_requests" // 确保与数据库表名一致
}

// BeforeSave 按工作小时计算（每日8小时）
// BeforeSave 在保存前自动计算时长（包括创建和更新）
func (lr *LeaveRequest) BeforeSave(tx *gorm.DB) (err error) {

	hours := lr.EndTime.Sub(lr.StartTime).Hours()
	// 每 8 小时为 1 天，保留一位小数（如 7.5 天）
	lr.Duration = math.Round(hours/8*10) / 10 // 四舍五入到一位小数
	return nil
}
