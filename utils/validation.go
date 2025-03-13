package utils

import (
	"github.com/go-playground/validator/v10"
)

func TranslateValidationErrors(err error) string {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			switch e.Tag() {
			case "required":
				return "字段 " + e.Field() + " 必须填写"
			case "min":
				return e.Field() + " 长度不能小于 " + e.Param()
			case "max":
				return e.Field() + " 长度不能超过 " + e.Param()
			case "email":
				return "邮箱格式不正确"
			case "len":
				return e.Field() + " 长度必须为 " + e.Param()
			}
		}
	}
	return "参数验证失败"
}
