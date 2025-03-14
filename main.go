package main

import (
	"EmployeeManagementDemo/config"
	"EmployeeManagementDemo/models"
	"EmployeeManagementDemo/routes"
	"EmployeeManagementDemo/services"
	"encoding/json"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// 初始化 MySQL 并自动迁移表结构
	setupDatabase()

	// 初始化 RabbitMQ
	config.InitRabbitMQ()
	defer config.RabbitMQConn.Close() // 程序退出时关闭连接
	defer config.RabbitMQChannel.Close()

	// 启动日志消费者
	services.StartLogConsumer()

	// 初始化 Gin 引擎
	router := gin.Default()

	// 消息队列测试
	//sendMessageTest()
	//
	//receiveMessageTest()

	// ---------- 新增：配置 CORS ----------
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:3000",       // 前端开发环境地址
			"https://your-production.com", // 生产环境地址（按需修改）
		},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Authorization",
			"X-Requested-With",
		},
		AllowCredentials: true,           // 允许携带 Cookie
		MaxAge:           12 * time.Hour, // 预检请求缓存时间
	}))
	// ---------- 新增结束 ----------

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
	setupRedis()
	//defer config.CloseMySQL()

	// 调整迁移顺序确保基础表先创建
	err := config.DB.AutoMigrate(
		&models.Department{},
		&models.Admin{},
		&models.Employee{},
		&models.SignRecord{},
		&models.LeaveRequest{},
		&models.OperationLog{},
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
	//defer config.Rdb.Close()

	if err := config.PingRedis(); err != nil {
		log.Fatalf("Redis 连接失败: %v", err)
	} else {
		log.Println("Redis 连接成功")
	}
}

func sendMessageTest() {
	// 发送测试消息
	message := map[string]string{"action": "test", "user_id": "123"}
	body, _ := json.Marshal(message)

	err := config.RabbitMQChannel.Publish(
		"",               // 使用默认交换机
		"operation_logs", // 队列名
		false,            // 强制标志
		false,            // 立即标志
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		log.Fatal("消息发送失败:", err)
	}

	log.Println("测试消息已发送")
}

func receiveMessageTest() {
	// 启动消费者
	msgs, err := config.RabbitMQChannel.Consume(
		"operation_logs", // 队列名
		"",               // 消费者标签
		true,             // 自动确认（消息处理后自动删除）
		false,            // 排他
		false,            // 无本地
		false,            // 无等待
		nil,
	)
	if err != nil {
		log.Fatal("注册消费者失败:", err)
	}

	// 异步处理消息
	go func() {
		for msg := range msgs {
			var data map[string]string
			if err := json.Unmarshal(msg.Body, &data); err != nil {
				log.Printf("消息解析失败: %v", err)
				continue
			}
			log.Printf("收到消息: %+v", data)
		}
	}()
}
