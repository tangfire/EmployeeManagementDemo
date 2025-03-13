package config

import (
	"context"
	"database/sql"
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

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

var Rdb *redis.Client
var Ctx = context.Background()

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
