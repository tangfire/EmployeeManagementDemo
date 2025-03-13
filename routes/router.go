// routes/router.go
package routes

import (
	"EmployeeManagementDemo/controllers"
	"EmployeeManagementDemo/middleware"
	"github.com/gin-gonic/gin"
)

// routes/router.go
func SetupAuthRoutes(r *gin.Engine) {
	//r.POST("/api/login", controllers.Login)
	//r.POST("/api/register", controllers.Register)            // 注册接口（员工自助）
	//r.POST("/api/admin/register", controllers.AdminRegister) // 管理员创建账号（需鉴权）

	// 公开接口（无需鉴权）
	publicGroup := r.Group("/api")
	{
		// 登录注册
		publicGroup.POST("/login", controllers.Login)
		publicGroup.POST("/register", controllers.Register)            // 员工自助注册
		publicGroup.POST("/admin/register", controllers.AdminRegister) // 管理员注册（需要密钥，但不需要登录）
	}

	// 需要登录的接口（需要鉴权）
	authGroup := r.Group("/api")
	authGroup.Use(middleware.JWTAuth())
	{
		// 员工个人信息
		authGroup.GET("/profile", controllers.GetProfile)
		//authGroup.PUT("/profile", controllers.UpdateProfile)

		// 签到/签退
		//authGroup.POST("/sign-records/sign-in", controllers.SignIn)
		//authGroup.POST("/sign-records/sign-out", controllers.SignOut)
	}

	// 需要管理员权限的接口（鉴权 + 管理员角色）
	adminGroup := r.Group("/api/admin")
	adminGroup.Use(middleware.JWTAuth(), middleware.AdminOnly())
	{
		//adminGroup.GET("/employees", controllers.ListEmployees)
		//adminGroup.POST("/departments", controllers.CreateDepartment)
	}

}
