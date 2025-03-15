package models

// models/user.go
type UserPasswordHolder interface {
	GetPassword() string
	SetPassword(string)
}

// 定义用户接口（统一Admin和Employee的认证行为）
type BaseUser interface {
	GetID() uint
	GetUsername() string
	GetRole() string
	CheckPassword(password string) bool
}
