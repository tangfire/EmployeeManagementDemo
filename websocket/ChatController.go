package websocket

import (
	"EmployeeManagementDemo/config"
	"EmployeeManagementDemo/models"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"strconv"
	"time"
)

// 聊天的后端调度逻辑
// 单聊
func SingleChat(c *models.Client, sendMsg *models.SendMsg) {
	// 使用更明确的键名
	chatKey := fmt.Sprintf("chat:%s", c.ID)
	r1, _ := config.Rdb.Get(context.Background(), chatKey).Result()
	//从redis中取出固定用户发给当前用户的消息
	id := CreateId(strconv.Itoa(c.RecipientID), strconv.Itoa(c.SendID))
	r2, _ := config.Rdb.Get(context.Background(), id).Result()
	//根据redis的结果去做未关注聊天次数限制
	if r2 >= "3" && r1 == "" {
		ResponseWebSocket(c.Socket, CodeLimiteTimes, "未相互关注，限制聊天次数")
		return
	} else {
		//将消息写入redis
		config.Rdb.Incr(context.Background(), c.ID)
		//设置消息的过期时间
		_, _ = config.Rdb.Expire(context.Background(), c.ID, time.Hour*24*30*3).Result()
	}
	fmt.Println(c.ID+"发送消息：", sendMsg.Content)
	//将消息广播出去
	models.Manager.Broadcast <- &models.Broadcast{
		Client:  c,
		Message: []byte(sendMsg.Content),
	}
}

// 查看未读消息
func UnreadMessages(c *models.Client) {
	//获取数据库中的未读消息
	msgs, err := GetMessageUnread(c.SendID)
	if err != nil {
		ResponseWebSocket(c.Socket, CodeServerBusy, "服务繁忙")
	}
	for i, msg := range msgs {
		replyMsg := models.ReplyMsg{
			From:    msg.Direction,
			Content: msg.Content,
		}
		message, _ := json.Marshal(replyMsg)
		_ = c.Socket.WriteMessage(websocket.TextMessage, message)
		//发送完后将消息设为已读
		msgs[i].Read = true
		err := UpdateMessage(&msgs[i])
		if err != nil {
			ResponseWebSocket(c.Socket, CodeServerBusy, "服务繁忙")
		}
	}
}

// 拉取历史消息记录
func HistoryMsg(c *models.Client, sendMsg *models.SendMsg) {
	//拿到传过来的时间
	timeT := TimeStringToGoTime(sendMsg.Content)
	//查找聊天记录
	//做一个分页处理，一次查询十条数据,根据时间去限制次数
	//别人发给当前用户的
	direction := CreateId(strconv.Itoa(c.RecipientID), strconv.Itoa(c.SendID))
	//当前用户发出的
	id := CreateId(strconv.Itoa(c.SendID), strconv.Itoa(c.RecipientID))
	msgs, err := GetHistoryMsg(direction, id, timeT, 10)
	if err != nil {
		ResponseWebSocket(c.Socket, CodeServerBusy, "服务繁忙")
	}
	//把消息写给用户
	for _, msg := range *msgs {
		replyMsg := models.ReplyMsg{
			From:    msg.Direction,
			Content: msg.Content,
		}
		message, _ := json.Marshal(replyMsg)
		_ = c.Socket.WriteMessage(websocket.TextMessage, message)
		//发送完后将消息设为已读
		if err != nil {
			ResponseWebSocket(c.Socket, CodeServerBusy, "服务繁忙")
		}
	}
}

// 群聊消息广播
func GroupChat(c *models.Client, sendMsg *models.SendMsg) {
	//根据消息类型判断是否为群聊消息
	//先去数据库查询该群下的所有用户
	users, err := GetAllGroupUser(strconv.Itoa(sendMsg.RecipientID))
	if err != nil {
		ResponseWebSocket(c.Socket, CodeServerBusy, "服务繁忙")
	}
	//向群里面的用户广播消息
	for _, user := range users {
		//获取群里每个用户的连接
		if int(user.ID) == c.SendID {
			continue
		}
		c.ID = strconv.Itoa(c.SendID) + "->"
		c.GroupID = strconv.Itoa(sendMsg.RecipientID)
		c.RecipientID = int(user.ID)
		models.Manager.Broadcast <- &models.Broadcast{
			Client:  c,
			Message: []byte(sendMsg.Content),
		}
	}
}
