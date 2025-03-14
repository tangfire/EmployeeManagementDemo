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

	userGroup := r.Group("/api")
	userGroup.Use(middleware.JWTAuth(), middleware.CheckJWTBlacklist())
	{
		userGroup.GET("/departments", controllers.GetDepartments)
		userGroup.PUT("/profile/password", controllers.UpdatePassword) // 新增密码修改路由
		// 个人信息
		userGroup.GET("/profile", controllers.GetProfile)
		userGroup.PUT("/profile", controllers.UpdateProfile)

		userGroup.POST("/logout", controllers.Logout) // 用户注销接口
	}

	// 需要登录的接口（需要鉴权）
	employeeGroup := r.Group("/api")
	employeeGroup.Use(middleware.JWTAuth(), middleware.CheckJWTBlacklist(), middleware.RequireRole("employee"))
	{

		// 签到签退路由（仅允许员工角色）
		//authGroup.POST("/sign-records/sign-in", middleware.RequireRole("employee"), controllers.SignIn)
		//authGroup.POST("/sign-records/sign-out", middleware.RequireRole("employee"), controllers.SignOut)
		employeeGroup.POST("/sign-records/sign-in", controllers.SignIn)
		employeeGroup.POST("/sign-records/sign-out", controllers.SignOut)

		employeeGroup.POST("/leave", controllers.CreateLeaveRequest) // 提交请假
		employeeGroup.GET("/leave", controllers.GetMyLeaveRequests)  // 查看自己的请假记录

	}

	// 需要管理员权限的接口（鉴权 + 管理员角色）
	adminGroup := r.Group("/api/admin")
	adminGroup.Use(middleware.JWTAuth(), middleware.CheckJWTBlacklist(), middleware.RequireRole("admin"))
	{
		//adminGroup.GET("/employees", controllers.ListEmployees)

		adminGroup.POST("/departments", controllers.CreateDepartment)           // 创建部门
		adminGroup.PUT("/departments/:dep_id", controllers.UpdateDepartment)    // 更新部门
		adminGroup.DELETE("/departments/:dep_id", controllers.DeleteDepartment) // 删除部门
		adminGroup.POST("/employees", controllers.CreateEmployee)               // POST   /api/employees
		adminGroup.PUT("/employees/:emp_id", controllers.UpdateEmployee)        // PUT    /api/employees/:emp_id
		adminGroup.DELETE("/employees/:emp_id", controllers.DeleteEmployee)     // DELETE /api/employees/:emp_id

		adminGroup.GET("/employees", controllers.GetEmployees) // GET    /api/employees

		adminGroup.PUT("/leave/:id/approve", controllers.ApproveLeaveRequest) // 审批
		adminGroup.GET("/leaves", controllers.GetAllLeaveRequests)            // 查看所有记录

		// 管理员踢人接口（需要管理员权限）
		adminGroup.PUT("/users/:user_id/kick", controllers.KickUser)

	}

}
