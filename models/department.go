// Package models models/department.go
package models

type Department struct {
	DepID  int    `gorm:"primaryKey;autoIncrement;column:dep_id"` // 主键自增
	Depart string `gorm:"type:varchar(20);not null"`              // 部门名称（非空）
}

// TableName 自定义表名（与数据库表名一致）
func (Department) TableName() string {
	return "departments"
}
