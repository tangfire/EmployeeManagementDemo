// Package models models/department.go
package models

type Department struct {
	DepID  uint   `gorm:"primaryKey;autoIncrement;column:dep_id"` // 主键自增
	Depart string `gorm:"type:varchar(20);not null;uniqueIndex"`  // 部门名称唯一
	// 移除 DeletedAt 字段（部门一般需要保留历史数据）
	Employees []Employee `gorm:"foreignKey:DepID;constraint:OnDelete:SET NULL;"`
}

// TableName 自定义表名（与数据库表名一致）
func (Department) TableName() string {
	return "departments"
}
