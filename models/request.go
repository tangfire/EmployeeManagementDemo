package models

import "time"

// controllers/LoginAuth.go
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=4,max=20"`
	Password string `json:"password" binding:"required,min=6,max=20"`
	Email    string `json:"email" binding:"required,email"`
	Phone    string `json:"phone" binding:"required,len=11"`
}

type AdminRegisterRequest struct {
	RegisterRequest
	SecretKey string `json:"secret_key" binding:"required"` // 必须携带密钥
}

// controllers/LoginAuth.go
type LoginRequest struct {
	Username string `json:"username" binding:"required,min=4,max=20"`
	Password string `json:"password" binding:"required,min=6,max=20"`
}

// models/request.go
type UpdateProfileRequest struct {
	Name   string `json:"username" binding:"required,min=4,max=20"` // 必填
	Email  string `json:"email" binding:"required,email"`
	Phone  string `json:"phone" binding:"omitempty,len=11"` // 可选，长度11
	Avatar string `json:"avatar" binding:"omitempty,url"`   // 头像URL
}

type UpdatePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`       // 必须提供旧密码
	NewPassword string `json:"new_password" binding:"required,min=6"` // 新密码至少6位
}

// 创建部门请求
type CreateDepartmentRequest struct {
	Depart string `json:"depart" binding:"required,min=1,max=20"` // 必填，长度1-20字符
}

// 更新部门请求
type UpdateDepartmentRequest struct {
	Depart string `json:"depart" binding:"required,min=1,max=20"`
}

// 创建员工请求
type CreateEmployeeRequest struct {
	Name         string `json:"username" binding:"required,min=2"` // 必填
	DepartmentID uint   `json:"department_id" binding:"required"`  // 必填
	Position     string `json:"position" binding:"required"`       // 必填
	Email        string `json:"email" binding:"omitempty,email"`   // 可选
	Phone        string `json:"phone" binding:"omitempty,len=11"`  // 可选
}

// 更新员工请求
type UpdateEmployeeRequest struct {
	Name         string `json:"username" binding:"omitempty,min=2"`
	DepartmentID uint   `json:"department_id" binding:"omitempty"`
	Position     string `json:"position"`
	Email        string `json:"email" binding:"omitempty,email"`
	Phone        string `json:"phone" binding:"omitempty,len=11"`
	Status       string `json:"status" binding:"omitempty,oneof=在职 离职"`
}

// 员工提交请假请求
type CreateLeaveRequest struct {
	Reason    string    `json:"reason" binding:"required,min=5"`               // 请假原因
	StartTime time.Time `json:"start_time" binding:"required"`                 // 开始时间
	EndTime   time.Time `json:"end_time" binding:"required,gtfield=StartTime"` // 结束时间需晚于开始时间
}

// 管理员审批请求
type ApproveLeaveRequest struct {
	Status  string `json:"status" binding:"required,oneof=approved rejected"` // 审批结果
	Comment string `json:"comment" binding:"omitempty,max=200"`               // 审批意见（可选）
}
