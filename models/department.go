// Package models models/department.go
package models

type Department struct {
	DepID  uint   `gorm:"primaryKey;autoIncrement;column:dep_id" json:"dep_id"` // 主键自增
	Depart string `gorm:"type:varchar(20);not null;uniqueIndex;" json:"depart"` // 部门名称唯一

}

// TableName 自定义表名（与数据库表名一致）
func (Department) TableName() string {
	return "departments"
}
