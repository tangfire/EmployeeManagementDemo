package config

import (
	"context"
	"database/sql"
	"github.com/go-redis/redis/v8"
	amqp "github.com/rabbitmq/amqp091-go"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

var (
	DB              *gorm.DB
	Rdb             *redis.Client
	Ctx             = context.Background()
	RabbitMQConn    *amqp.Connection // RabbitMQ 连接对象
	RabbitMQChannel *amqp.Channel    // RabbitMQ 通道对象
)

// 初始化 RabbitMQ 连接
func InitRabbitMQ() {
	var err error

	// 1. 连接到 RabbitMQ 服务器
	RabbitMQConn, err = amqp.Dial("amqp://admin:8888.216@localhost:5674/my_vhost")
	if err != nil {
		log.Fatal("RabbitMQ 连接失败:", err)
	}

	// 2. 创建通道
	RabbitMQChannel, err = RabbitMQConn.Channel()
	if err != nil {
		log.Fatal("创建 RabbitMQ 通道失败:", err)
	}

	// 3. 声明队列（如果队列不存在则创建）
	_, err = RabbitMQChannel.QueueDeclare(
		"operation_logs", // 队列名称
		true,             // 持久化（重启后队列仍存在）
		false,            // 自动删除（无消费者时自动删除）
		false,            // 排他队列（仅当前连接可见）
		false,            // 无等待（不等待服务器响应）
		nil,              // 额外参数
	)
	if err != nil {
		log.Fatal("声明队列失败:", err)
	}

	log.Println("RabbitMQ 初始化成功")
}

func InitMySQL() {
	dsn := "root:8888.216@tcp(localhost:3306)/employee_db?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		// 禁用自动创建外键约束
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		panic("Failed to connect to MySQL: " + err.Error())
	}
	DB = db
}

// GetSQLDB 新增获取底层 *sql.DB 的方法
func GetSQLDB() *sql.DB {
	sqlDB, err := DB.DB()
	if err != nil {
		panic("获取 SQL DB 失败: " + err.Error())
	}
	return sqlDB
}

// CloseMySQL 新增关闭连接的方法
func CloseMySQL() {
	if sqlDB := GetSQLDB(); sqlDB != nil {
		sqlDB.Close()
	}
}

// PingMySQL config/config.go 新增
func PingMySQL() error {
	return GetSQLDB().Ping()
}

func InitRedis() {
	Rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "8888.216",
		DB:       8,
	})
}

// PingRedis utils/redis.go 新增
func PingRedis() error {
	return Rdb.Ping(Ctx).Err()
}
