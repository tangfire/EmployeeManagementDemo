package models

import (
	"time"
)

type SignRecord struct {
	ID          uint       `gorm:"primaryKey;autoIncrement;column:id"`
	EmpID       uint       `gorm:"uniqueIndex:idx_emp_date"`
	SignInTime  *time.Time `gorm:"type:datetime;"`
	SignOutTime *time.Time `gorm:"type:datetime;"`
	Date        time.Time  `gorm:"uniqueIndex:idx_emp_date;type:date"`
	Status      string     `gorm:"type:enum('正常','迟到','早退','缺勤');default:'缺勤'"`
}

func (SignRecord) TableName() string {
	return "sign_records"
}
