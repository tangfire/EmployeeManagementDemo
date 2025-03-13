package models

import "golang.org/x/crypto/bcrypt"

type Admin struct {
	AdminID       uint   `gorm:"primaryKey;autoIncrement;column:admin_id"` // 主键自增
	AdminName     string `gorm:"type:varchar(20);not null;unique"`         // 不能为空且唯一
	AdminPassword string `gorm:"type:varchar(200);not null"`               // 不能为空
	AdminEmail    string `gorm:"type:varchar(50)"`                         // 邮箱
	AdminPhone    string `gorm:"type:varchar(11)"`                         // 手机号
	Avatar        string `gorm:"type:varchar(100)"`                        // 头像路径

	// 可选：添加与 LeaveRequest 的反向关联
	LeaveRequests []LeaveRequest `gorm:"foreignKey:AdminID"`
}

// TableName 自定义表名（与数据库表名一致）
func (Admin) TableName() string {
	return "admin"
}

func (a *Admin) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(a.AdminPassword), []byte(password))
	return err == nil
}

// models/admin.go
func (a *Admin) GetID() uint {
	return a.AdminID
}

func (a *Admin) GetUsername() string {
	return a.AdminName
}

func (a *Admin) GetRole() string {
	return "admin" // 标识角色为管理员
}

type AdminProfileResponse struct {
	AdminID   uint   `json:"admin_id"`
	AdminName string `json:"username"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
}

func (a *Admin) ToProfileResponse() AdminProfileResponse {
	return AdminProfileResponse{
		AdminID:   a.AdminID,
		AdminName: a.AdminName,
		Email:     a.AdminEmail,
		Phone:     a.AdminPhone,
	}
}
