// websocket/manager.go
package websocket

import (
	"EmployeeManagementDemo/utils"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"sync"
	"time"
)

// websocket/manager.go
var onlineUsers = make(map[uint]bool)
var onlineMutex sync.RWMutex

// 获取在线用户列表
func GetOnlineUsers() map[uint]bool {
	onlineMutex.RLock()
	defer onlineMutex.RUnlock()
	return onlineUsers
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Client struct {
	UserID uint
	Conn   *websocket.Conn
}

var (
	clients   = make(map[uint]*Client)
	clientsMu sync.RWMutex
)

func HandleWebSocket(c *gin.Context) {
	// 用户鉴权
	userID, err := utils.GetCurrentUserID(c)
	if err != nil {
		c.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	// 升级WebSocket连接
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": "WebSocket upgrade failed"})
		return
	}

	// 在HandleWebSocket连接成功后添加
	onlineMutex.Lock()
	onlineUsers[userID] = true
	onlineMutex.Unlock()

	// 注册客户端
	client := &Client{UserID: userID, Conn: conn}
	clientsMu.Lock()
	clients[userID] = client
	clientsMu.Unlock()

	// 消息处理循环
	go handleClient(client)

	// 在断开连接时删除
	defer func() {
		onlineMutex.Lock()
		delete(onlineUsers, userID)
		onlineMutex.Unlock()
	}()

}

func handleClient(client *Client) {
	defer func() {
		clientsMu.Lock()
		delete(clients, client.UserID)
		clientsMu.Unlock()
		client.Conn.Close()
	}()

	for {
		_, msgBytes, err := client.Conn.ReadMessage()
		if err != nil {
			break
		}

		var message struct {
			ReceiverID uint   `json:"receiverId"`
			Content    string `json:"content"`
			Timestamp  int64  `json:"timestamp"`
		}

		if err := json.Unmarshal(msgBytes, &message); err != nil {
			continue
		}

		// 添加时间戳
		message.Timestamp = time.Now().Unix()

		// 发送给接收者
		sendMessage(message.ReceiverID, message)
		// 同时发送给发送者（实现消息同步）
		sendMessage(client.UserID, message)
	}
}

func sendMessage(receiverID uint, message interface{}) {
	clientsMu.RLock()
	defer clientsMu.RUnlock()

	if client, ok := clients[receiverID]; ok {
		msgBytes, _ := json.Marshal(message)
		client.Conn.WriteMessage(websocket.TextMessage, msgBytes)
	}
}
