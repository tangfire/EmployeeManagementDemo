package services

import (
	"EmployeeManagementDemo/config"
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

// 通用消息发送函数
func SendLogToRabbitMQ(data map[string]interface{}) {
	body, _ := json.Marshal(data)
	err := config.RabbitMQChannel.Publish(
		"",               // 使用默认交换机
		"operation_logs", // 队列名
		false,            // 强制标志
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent, // 持久化消息
		},
	)
	if err != nil {
		log.Printf("日志发送失败: %v", err)
	}
}
