package utils

import (
	"EmployeeManagementDemo/services"
	"github.com/dgrijalva/jwt-go"
	"time"
)

var JwtSecret = []byte("your-secret-key") // 从环境变量读取更安全

type Claims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	jwt.StandardClaims
}

func GenerateJWT(user services.User) (string, error) {
	claims := Claims{
		UserID: user.GetID(),
		Role:   user.GetRole(),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
			Issuer:    "EmployeeManagement",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JwtSecret)
}
