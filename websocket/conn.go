// websocket/conn.go
package websocket

import (
	"EmployeeManagementDemo/models"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"sync"
)

// 在全局或结构体中声明
var rwLocker sync.RWMutex

func WsHandle(c *gin.Context) {
	myid := c.Query("myid")
	userid, err := strconv.Atoi(myid)
	if err != nil {
		zap.L().Error("转换失败", zap.Error(err))
		ResponseError(c, CodeParamError)
	}
	//将http协议升级为ws协议
	conn, err := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
		return true
	}}).Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		http.NotFound(c.Writer, c.Request)
		return
	}
	//创建一个用户客户端实例，用于记录该用户的连接信息
	client := new(models.Client)
	client = &models.Client{
		ID:     myid + "->",
		SendID: userid,
		Socket: conn,
		Send:   make(chan []byte),
	}
	//使用管道将实例注册到用户管理上
	models.Manager.Register <- client
	//开启两个协程用于读写消息
	go Read(client)
	go Write(client)
}

// 用于读管道中的数据
func Read(c *models.Client) {
	//结束把通道关闭
	defer func() {
		models.Manager.Unregister <- c
		//关闭连接
		_ = c.Socket.Close()
	}()
	for {
		//先测试一下连接能不能连上
		c.Socket.PongHandler()
		sendMsg := new(models.SendMsg)
		err := c.Socket.ReadJSON(sendMsg)
		c.RecipientID = sendMsg.RecipientID
		if err != nil {
			zap.L().Error("数据格式不正确", zap.Error(err))
			models.Manager.Unregister <- c
			_ = c.Socket.Close()
			return
		}
		//根据要发送的消息类型去判断怎么处理
		//消息类型的后端调度
		switch sendMsg.Type {
		case 1: //私信
			SingleChat(c, sendMsg)
		case 2: //获取未读消息
			UnreadMessages(c)
		case 3: //拉取历史消息记录
			HistoryMsg(c, sendMsg)
		case 4: //群聊消息广播
			GroupChat(c, sendMsg)
		}
	}
}

// websocket/conn.go
func Write(c *models.Client) {
	defer func() {
		_ = c.Socket.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				_ = c.Socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			replyMsg := models.ReplyMsg{
				Code:    CodeConnectionSuccess,
				Content: string(message),
			}
			msg, _ := json.Marshal(replyMsg)

			rwLocker.Lock()
			_ = c.Socket.WriteMessage(websocket.TextMessage, msg)
			rwLocker.Unlock() // 确保解锁
		}
	}
}
