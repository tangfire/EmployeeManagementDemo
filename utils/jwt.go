package utils

import (
	"EmployeeManagementDemo/models"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

var JwtSecret = []byte("your-secret-key")

type Claims struct {
	UserID               uint   `json:"user_id"`
	Role                 string `json:"role"`
	jwt.RegisteredClaims        // 替换 StandardClaims
}

func GenerateJWT(user models.User) (string, error) {
	claims := Claims{
		UserID: user.GetID(),
		Role:   user.GetRole(),
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()), // 新增：设置签发时间
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			Issuer:    "EmployeeManagement",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JwtSecret)
}

func ParseJWT(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return JwtSecret, nil
	})

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}
