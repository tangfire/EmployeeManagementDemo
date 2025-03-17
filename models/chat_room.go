package models

import (
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
	"time"
)

// ChatMessage 数据库存储消息结构体，用于持久化历史记录
type ChatMessage struct {
	gorm.Model
	Direction   string //这条消息是从谁发给谁的
	SendID      int    //发送者id
	RecipientID int    //接受者id
	GroupID     string //群id，该消息要发到哪个群里面去
	Content     string //内容
	Read        bool   //是否读了这条消息
}

// Group 群聊结构体
type Group struct {
	ID           string ` gorm:"primaryKey"` //群id
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
	GroupName    string         `json:"group_name"`    //群名
	GroupContent string         `json:"group_content"` //群签名
	GroupIcon    string         `json:"group_icon"`    //群头像
	GroupNum     int            //群人数
	GroupOwnerId int            //群主id
	Users        []User         `gorm:"many2many:users_groups;"` //群成员
}

type UsersGroup struct {
	GroupId string `json:"group_id"`
	UserId  int    `json:"user_id"`
}

// ReplyMsg 用于处理请求后返回一些数据
type ReplyMsg struct {
	From    string `json:"from"`
	Code    int    `json:"code"`
	Content string `json:"content"`
}

// SendMsg 发送消息的类型
type SendMsg struct {
	Type        int    `json:"type"`
	RecipientID int    `json:"recipient_id"` //接受者id
	Content     string `json:"content"`
}

// Client 用户类
type Client struct {
	ID          string          //消息的去向
	RecipientID int             //接受者id
	SendID      int             //发送人的id
	GroupID     string          //群聊id
	Socket      *websocket.Conn //websocket连接对象
	Send        chan []byte     //发送消息用的管道
}

// Broadcast 广播类，包括广播内容和源用户
type Broadcast struct {
	Client  *Client
	Message []byte
	Type    int
}

// ClientManager 用户管理,用于管理用户的连接及断开连接
type ClientManager struct {
	Clients    map[string]*Client
	Broadcast  chan *Broadcast
	Reply      chan *Client
	Register   chan *Client
	Unregister chan *Client
}

// Manager 创建一个用户管理对象
var Manager = ClientManager{
	Clients:    make(map[string]*Client), // 参与连接的用户，出于性能的考虑，需要设置最大连接数
	Broadcast:  make(chan *Broadcast),
	Register:   make(chan *Client), //新建立的连接访放入这里面
	Reply:      make(chan *Client),
	Unregister: make(chan *Client), //新断开的连接放入这里面
}
