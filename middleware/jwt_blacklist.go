// middleware/jwt_blacklist.go
package middleware

import (
	"EmployeeManagementDemo/config"
	"EmployeeManagementDemo/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
)

// middleware/jwt_blacklist.go
func CheckJWTBlacklist() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 检查 Token 是否在黑名单
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "未提供令牌"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// 检查 Token 黑名单
		exists, err := config.Rdb.Exists(config.Ctx, "jwt_blacklist:"+tokenString).Result()
		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": "服务器错误"})
			return
		}
		if exists == 1 {
			c.AbortWithStatusJSON(401, gin.H{"error": "令牌已失效"})
			return
		}

		// 2. 检查用户是否被踢出（需依赖 JWTAuth 设置的 claims）
		claimsInterface, ex := c.Get("claims")
		if !ex {
			c.AbortWithStatusJSON(401, gin.H{"error": "令牌无效"})
			return
		}

		claims, ok := claimsInterface.(*utils.Claims)
		if !ok {
			c.AbortWithStatusJSON(500, gin.H{"error": "服务器内部错误"})
			return
		}

		// 检查用户是否被踢出
		userID := claims.UserID
		invalidateKey := "user_invalid:" + fmt.Sprint(userID)
		kickTime, err := config.Rdb.Get(config.Ctx, invalidateKey).Int64()
		if err == nil { // 存在踢出记录
			tokenIssueTime := claims.IssuedAt.Unix()
			if tokenIssueTime < kickTime {
				c.AbortWithStatusJSON(401, gin.H{"error": "用户已被踢出"})
				return
			}
		}

		c.Next()
	}
}
