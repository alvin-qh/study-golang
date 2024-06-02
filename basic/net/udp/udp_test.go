package udp

/**
 * Golang with UDP
 *
 * UDP 连接和 TCP 连接有如下区别:
 *  1. 无状态, 无客户端连接, 所有的数据接收和发送均通过一个 UDP 连接处理
 *  2. 接收和发送数据需按照顺序
 * 服务端处理过程如下:
 *  1. 启动监听, 得到一个 UDP 连接对象
 *  2. goroutine 1: 在 UDP 连接对象上启动数据接收
 *  3. goroutine 2: 监听一个发送 channel, 并将 channel 发来的数据通过 UDP 发送给客户端
 *  4. 接收的数据报带有发送方地址, 服务端需对每个发送方地址做标识, 以记录其状态
 */

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func TestUDP_Network(t *testing.T) {
	// 启动服务器
	server, err := ServerStart("0.0.0.0:8888")
	assert.Nil(t, err)
	defer server.Close()

	// 连接服务器
	client, err := Connect("127.0.0.1:8888")
	assert.Nil(t, err)
	defer client.Close()

	// 发送登录请求
	resp, err := client.Request(&LoginAsk{
		AskHeader: AskHeader{Action: ACTION_LOGIN},
		Account:   "Alvin",
		Password:  "password",
	})
	assert.Nil(t, err)

	// 接收登录响应
	loginAck, ok := resp.(*LoginAck)
	assert.True(t, ok)
	assert.Equal(t, "Welcome Alvin", loginAck.Welcome)

	// 发送关闭服务请求
	resp, err = client.Request(&ShutdownAsk{
		AskHeader: AskHeader{Action: ACTION_SHUTDOWN},
	})
	assert.Nil(t, err)
	assert.Equal(t, ACTION_SHUTDOWN, resp.GetAction())
}
