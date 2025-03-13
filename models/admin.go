// Package models models/admin.go
package models

type Admin struct {
	AdminID       uint   `gorm:"primaryKey;autoIncrement;column:admin_id"` // 主键自增
	AdminName     string `gorm:"type:varchar(20);not null;unique"`         // 不能为空且唯一
	AdminPassword string `gorm:"type:varchar(20);not null"`                // 不能为空
	AdminEmail    string `gorm:"type:varchar(50)"`                         // 邮箱
	AdminPhone    string `gorm:"type:varchar(11)"`                         // 手机号
	Avatar        string `gorm:"type:varchar(100)"`                        // 头像路径
}

// TableName 自定义表名（与数据库表名一致）
func (Admin) TableName() string {
	return "admin"
}
