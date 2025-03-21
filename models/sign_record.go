package models

import (
	"database/sql"
	"time"
)

type SignRecord struct {
	ID          uint         `gorm:"primaryKey;autoIncrement;column:id"`
	EmpID       uint         `gorm:"uniqueIndex:idx_emp_date"`
	SignInTime  time.Time    `gorm:"type:datetime;default:CURRENT_TIMESTAMP"`
	SignOutTime sql.NullTime `gorm:"type:datetime;default:NULL"` // 允许 NULL
	Date        time.Time    `gorm:"uniqueIndex:idx_emp_date;type:date"`
	Status      string       `gorm:"type:enum('正常','迟到','早退','迟到+早退','缺勤');default:'缺勤'"`
}

func (SignRecord) TableName() string {
	return "sign_records"
}
