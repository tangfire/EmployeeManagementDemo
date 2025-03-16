package config

import (
	"context"
	"database/sql"
	"github.com/go-redis/redis/v8"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

type Config struct {
	App      AppConfig      `mapstructure:"app"`
	Database DatabaseConfig `mapstructure:"database"`
	RabbitMQ RabbitMQConfig `mapstructure:"rabbitmq"`
	Logging  LoggingConfig  `mapstructure:"logging"`
}

type AppConfig struct {
	Env  string `mapstructure:"env"`
	Port int    `mapstructure:"port"`
}

type MySQLConfig struct {
	DSN string `mapstructure:"dsn"`
}

type RedisConfig struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type DatabaseConfig struct {
	MySQL MySQLConfig `mapstructure:"mysql"`
	Redis RedisConfig `mapstructure:"redis"`
}

type RabbitMQConfig struct {
	URL   string `mapstructure:"url"`
	Queue string `mapstructure:"queue"`
}

type LoggingConfig struct {
	Level string `mapstructure:"level"`
}

var Cfg Config

func LoadConfig() {
	v := viper.New()

	// 基础配置
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("./config")
	v.AddConfigPath(".") // 兼容不同执行路径

	// 环境变量支持（优先级高于配置文件）
	v.AutomaticEnv()
	v.SetEnvPrefix("APP") // 环境变量前缀 APP_DATABASE_MYSQL_DSN

	// 读取配置
	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("读取配置文件失败: %v", err)
	}

	// 反序列化到结构体
	if err := v.Unmarshal(&Cfg); err != nil {
		log.Fatalf("配置解析失败: %v", err)
	}
}

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

	//1. 连接到 RabbitMQ 服务器
	//RabbitMQConn, err = amqp.Dial("amqp://admin:8888.216@localhost:5674/my_vhost")
	//if err != nil {
	//	log.Fatal("RabbitMQ 连接失败:", err)
	//}

	RabbitMQConn, err = amqp.Dial(Cfg.RabbitMQ.URL)
	if err != nil {
		log.Fatal("RabbitMQ连接失败:", err)
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
	//dsn := "root:8888.216@tcp(localhost:3306)/employee_db?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := Cfg.Database.MySQL.DSN
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		// 禁用自动创建外键约束
		//DisableForeignKeyConstraintWhenMigrating: true,
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
	//Rdb = redis.NewClient(&redis.Options{
	//	Addr:     "localhost:6379",
	//	Password: "8888.216",
	//	DB:       8,
	//})
	Rdb = redis.NewClient(&redis.Options{
		Addr:     Cfg.Database.Redis.Addr,
		Password: Cfg.Database.Redis.Password,
		DB:       Cfg.Database.Redis.DB,
	})
}

// PingRedis utils/redis.go 新增
func PingRedis() error {
	return Rdb.Ping(Ctx).Err()
}
