package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func RequireRole(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, ok := c.Get("userRole")
		if !ok || role.(string) != requiredRole {
			c.JSON(http.StatusForbidden, gin.H{"error": "无权限访问"})
			c.Abort()
			return
		}
		c.Next()
	}
}
