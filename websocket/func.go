package websocket

import (
	"EmployeeManagementDemo/config"
	"EmployeeManagementDemo/models"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"time"
)

func CreateId(id1, id2 string) string {
	// 按数值大小排序，确保键名唯一性（如 1001:2002 和 2002:1001 生成相同键）
	if id1 > id2 {
		id1, id2 = id2, id1
	}
	return fmt.Sprintf("%s:%s", id1, id2)
}

// 修正后的响应函数
func ResponseWebSocket(conn *websocket.Conn, code int, message string) error {
	// 验证连接有效性
	if conn == nil {
		return errors.New("websocket连接未初始化")
	}

	response := struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}{
		Code:    code,
		Message: message,
	}

	// 必须处理JSON序列化错误
	data, err := json.Marshal(response)
	if err != nil {
		return fmt.Errorf("响应序列化失败: %w", err)
	}

	// 使用线程安全写入
	return conn.WriteMessage(websocket.TextMessage, data)
}

func GetMessageUnread(recipientID int) ([]models.ChatMessage, error) {
	var messages []models.ChatMessage

	// 使用 GORM 查询构建器
	result := config.DB.Where(
		"recipient_id = ? AND `read` = ?",
		recipientID,
		false,
	).Order("created_at desc").Find(&messages)

	if result.Error != nil {
		return nil, result.Error
	}
	return messages, nil
}

// mysql/message.go
func UpdateMessage(msg *models.ChatMessage) error {
	// 使用Select限定仅更新Read字段
	result := config.DB.Model(msg).
		Select("Read").
		Where("id = ?", msg.ID).
		Updates(map[string]interface{}{"read": true})

	if result.Error != nil {
		return fmt.Errorf("消息状态更新失败: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func TimeStringToGoTime(timeStr string) time.Time {
	// 尝试解析RFC3339格式（如"2025-03-17T15:04:05Z"）
	if t, err := time.Parse(time.RFC3339, timeStr); err == nil {
		return t.UTC()
	}

	// 尝试解析UNIX时间戳（如"1740215045"）
	if unix, err := strconv.ParseInt(timeStr, 10, 64); err == nil {
		return time.Unix(unix, 0).UTC()
	}

	// 默认返回当前时间（根据网页3建议使用UTC）
	return time.Now().UTC()
}

func GetHistoryMsg(direction, id string, before time.Time, limit int) (*[]models.ChatMessage, error) {
	var messages []models.ChatMessage
	result := config.DB.
		Where("(direction = ? OR direction = ?) AND created_at < ?", direction, id, before).
		Order("created_at DESC").
		Limit(limit).
		Find(&messages)

	if result.Error != nil {
		return nil, result.Error
	}
	return &messages, nil
}

// websocket/func.go
func GetAllGroupUser(groupID string) ([]models.User, error) {
	var group models.Group
	result := config.DB.Preload("Users").Where("id = ?", groupID).First(&group)
	if result.Error != nil {
		return nil, fmt.Errorf("群组查询失败: %w", result.Error)
	}
	return group.Users, nil
}

func ResponseError(c *gin.Context, code int) {
	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
		"code":    code,
		"message": getErrorMessage(code),
	})
}

func getErrorMessage(code int) string {
	switch code {
	case CodeParamError:
		return "参数格式错误"
	default:
		return "未知错误"
	}
}
