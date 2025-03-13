// routes/router.go
package routes

import (
	"EmployeeManagementDemo/controllers"
	"github.com/gin-gonic/gin"
)

// routes/router.go
func SetupAuthRoutes(r *gin.Engine) {
	r.POST("/api/login", controllers.Login)
	r.POST("/api/register", controllers.Register)            // 注册接口（员工自助）
	r.POST("/api/admin/register", controllers.AdminRegister) // 管理员创建账号（需鉴权）
}
