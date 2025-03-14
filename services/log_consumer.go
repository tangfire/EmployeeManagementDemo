// services/log_consumer.go
package services

import (
	"EmployeeManagementDemo/config"
	"EmployeeManagementDemo/models"
	"encoding/json"
	"log"
	"time"
)

func StartLogConsumer() {
	// 声明队列（确保存在）
	_, err := config.RabbitMQChannel.QueueDeclare(
		"operation_logs", // 队列名
		true,             // 持久化
		false,            // 自动删除
		false,            // 排他
		false,            // 无等待
		nil,
	)
	if err != nil {
		log.Fatal("队列声明失败:", err)
	}

	// 消费消息
	msgs, err := config.RabbitMQChannel.Consume(
		"operation_logs",
		"",
		false, // 关闭自动确认（手动确认）
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal("注册消费者失败:", err)
	}

	// 异步处理消息
	go func() {
		for msg := range msgs {
			var logData map[string]interface{}
			if err := json.Unmarshal(msg.Body, &logData); err != nil {
				log.Printf("日志解析失败: %v", err)
				msg.Nack(false, true) // 拒绝并重新入队
				continue
			}

			// 写入数据库
			logEntry := models.OperationLog{
				UserID:    uint(logData["user_id"].(float64)), // 注意类型断言
				Action:    logData["action"].(string),
				TargetID:  logData["target_id"].(string),
				CreatedAt: time.Now(),
			}
			if err := config.DB.Create(&logEntry).Error; err != nil {
				log.Printf("日志写入失败: %v", err)
				msg.Nack(false, true)
				continue
			}

			msg.Ack(false) // 确认消息已处理
		}
	}()
}
