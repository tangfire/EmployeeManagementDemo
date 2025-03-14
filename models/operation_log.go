// models/operation_log.go
package models

import "time"

type OperationLog struct {
	ID        uint   `gorm:"primaryKey"`
	UserID    uint   `gorm:"index"`    // 操作人ID（如果是登录/注销，即用户自己）
	Action    string `gorm:"size:100"` // 操作类型：login, logout, kick_user
	TargetID  string `gorm:"size:100"` // 被操作对象ID（如被踢用户ID，登录时留空）
	CreatedAt time.Time
}

// TableName
func (OperationLog) TableName() string {
	return "operation_logs" // 确保与数据库表名一致
}
