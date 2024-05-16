package tcp

import "errors"

// 定义业务代码
type ActionCode int

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

const (
	ACTION_LOGIN    ActionCode = iota // 登录业务码
	ACTION_SHUTDOWN                   // 关闭服务器业务码
)

// 请求头结构体
type AskHeader struct {
	Action ActionCode // 业务码
}

// 响应头结构体
type AckHeader struct {
	Action ActionCode // 业务码
	IsOk   bool       // 是否成功
	Error  string     // 错误信息
}

// 登录请求结构体
type LoginAsk struct {
	Account  string
	Password string
}

// 登陆响应结构体
type LoginAck struct {
	Welcome string
}

// 关闭请求结构体
type ShutdownAsk struct {
}

// 关闭响应结构体
type ShutdownAck struct {
}
