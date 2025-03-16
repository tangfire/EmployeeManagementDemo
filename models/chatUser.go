// models/employee.go
package models

type ChatUser struct {
	ID     uint   `json:"id" gorm:"primaryKey"`
	Name   string `json:"name"`
	Online bool   `json:"online" gorm:"-"`
}
