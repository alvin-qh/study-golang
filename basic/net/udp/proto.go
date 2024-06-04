package udp

import "errors"

// 定义业务代码
type ActionCode int

const (
	PACKAGE_LIMIT = 1024 * 60
)

// 业务代码转字符串
func (a ActionCode) String() string {
	switch a {
	case ACTION_LOGIN:
		return "ACTION_LOGIN"
	case ACTION_SHUTDOWN:
		return "ACTION_SHUTDOWN"
	default:
		panic(errors.New("invalid action code"))
	}
}

// 业务代码
const (
	ACTION_LOGIN    ActionCode = iota // 登录业务码
	ACTION_SHUTDOWN                   // 关闭服务器业务码
)

// SessionId 类型
type SessionId string

// 定义 Package 接口
type Package interface {
	GetAction() ActionCode            // 获取业务代码
	SetSessionId(sessionId SessionId) // 设置 SessionId
}

// 请求头结构体
type AskHeader struct {
	Action    ActionCode // 业务码
	SessionId SessionId  // 会话ID
}

// 获取请求数据的 Action
func (h *AskHeader) GetAction() ActionCode {
	return h.Action
}

// 设置请求数据的 SessionID
func (h *AskHeader) SetSessionId(sessionId SessionId) {
	h.SessionId = sessionId
}

// 响应头结构体
type AckHeader struct {
	Action    ActionCode // 业务码
	SessionId SessionId  // 会话ID
	IsOk      bool       // 是否成功
	Error     string     // 错误信息
}

// 获取响应数据的 Action
func (h *AckHeader) GetAction() ActionCode {
	return h.Action
}

// 设置响应数据的 SessionId
func (h *AckHeader) SetSessionId(sessionId SessionId) {
	h.SessionId = sessionId
}

// 登录请求结构体
type LoginAsk struct {
	AskHeader
	Account  string
	Password string
}

// 登陆响应结构体
type LoginAck struct {
	AckHeader
	Welcome string
}

// 关闭请求结构体
type ShutdownAsk struct {
	AskHeader
}

// 关闭响应结构体
type ShutdownAck struct {
	AckHeader
}
