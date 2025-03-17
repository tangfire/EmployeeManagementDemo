package models

import "time"

// ApiResponse 后端通用响应结构（utils/response.go）
type ApiResponse struct {
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data"`
	Timestamp int64       `json:"timestamp"`
}

func Success(data interface{}) ApiResponse {
	return ApiResponse{
		Code:      200,
		Message:   "success",
		Data:      data,
		Timestamp: time.Now().Unix(),
	}
}

func Error(code int, message string) ApiResponse {
	return ApiResponse{
		Code:      code,
		Message:   message,
		Timestamp: time.Now().Unix(),
	}
}

// EmployeeProfileResponse models/employee.go
type EmployeeProfileResponse struct {
	EmpID    uint   `json:"emp_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	DepID    uint   `json:"dep_id"`
}

// User 定义用户实体
type User struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
}

// LoginDTO 定义登录响应DTO
type LoginDTO struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

// models/department.go
type DepartmentAvgSalaryDTO struct {
	DepID      uint    `gorm:"column:dep_id" json:"dep_id"`
	Department string  `gorm:"column:depart" json:"depart"` // 明确映射列名
	AvgSalary  float64 `gorm:"column:avg_salary" json:"avg_salary"`
}

// models/department.go
type DepartmentHeadcountDTO struct {
	DepID      uint    `gorm:"column:dep_id" json:"dep_id"`
	Department string  `gorm:"column:depart" json:"depart"`
	Headcount  int     `gorm:"column:headcount" json:"headcount"`
	Percentage float64 `json:"percentage"` // 新增比例字段
}

type EmployeeWithDepNameDTO struct {
	Employee
	DepName string `json:"dep_name"` // 仅用于接收联表查询结果
}
