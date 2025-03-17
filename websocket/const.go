package websocket

const (
	CodeLimiteTimes       = 4003 // 未关注用户聊天次数限制
	CodeParamError        = 4000
	CodeServerBusy        = 503 // 使用标准HTTP状态码
	CodeConnectionSuccess = 200
	CodeConnectionBreak   = 4004 // 连接中断（客户端主动断开或网络问题）
)
