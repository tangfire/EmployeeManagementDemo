package main

import (
	"EmployeeManagementDemo/config"
	"EmployeeManagementDemo/models"
	"EmployeeManagementDemo/routes"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// 初始化 MySQL 并自动迁移表结构
	setupDatabase()

	// 初始化 Gin 引擎
	router := gin.Default()

	// 注册路由
	setupRoutes(router)

	// 启动 HTTP 服务器
	go startServer(router)

	// 优雅关机处理
	waitForShutdown()

	// 初始化 Redis（按需启用）
	// setupRedis()

	// 运行测试数据初始化（按需启用）
	// runTestData()
}

// 封装数据库初始化与迁移
func setupDatabase() {
	config.InitMySQL()
	//defer config.CloseMySQL()

	// 调整迁移顺序确保基础表先创建
	err := config.DB.AutoMigrate(
		&models.Department{},
		&models.Admin{},
		&models.Employee{},
		&models.SignRecord{},
		&models.LeaveRequest{},
		// 依赖员工和管理员

	)
	if err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	} else {
		log.Println("所有表已创建/更新")
	}
}

// 注册路由
func setupRoutes(router *gin.Engine) {
	// 认证相关路由
	routes.SetupAuthRoutes(router)

	// 其他模块路由（后续扩展）
	// routes.SetupEmployeeRoutes(router)
	// routes.SetupDepartmentRoutes(router)
}

// 启动 HTTP 服务器
func startServer(router *gin.Engine) {
	port := ":8080" // 默认端口
	log.Printf("服务启动中，监听端口 %s", port)
	if err := router.Run(port); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}

// 优雅关机处理
func waitForShutdown() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("接收到关机信号，服务正在关闭...")
	// 这里可以添加资源释放逻辑（如关闭数据库连接）
	os.Exit(0)
}

// 封装 Redis 初始化
func setupRedis() {
	config.InitRedis()
	defer config.Rdb.Close()

	if err := config.PingRedis(); err != nil {
		log.Fatalf("Redis 连接失败: %v", err)
	} else {
		log.Println("Redis 连接成功")
	}
}
