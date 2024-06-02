package tcp

/**
 * Golang with TCP
 *
 * Golang 处理 TCP 连接较为简单. 通过把每个客户端连接放到一个 goroutine 中处理, 简化了异步处理
 * 服务端处理过程如下:
 *  1. 启动监听 (Listening)
 *  2. 等待客户端连接 (Accepting)
 *  3. 为客户端连接启动 goroutine
 *  4. 通过客户端连接接收数据
 *  5. 通过客户端连接发送数据
 */

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func TestTCP_Network(t *testing.T) {
	// 启动服务器
	server, err := ServerStart("0.0.0.0:8888")
	assert.Nil(t, err)
	defer server.Close()

	// 连接服务器
	client, err := Connect("127.0.0.1:8888")
	assert.Nil(t, err)
	defer client.Close()

	// 发送登录请求
	resp, err := client.Request(ACTION_LOGIN, &LoginAsk{
		Account:  "Alvin",
		Password: "password",
	})
	assert.Nil(t, err)

	// 接收登录响应
	loginAck, ok := resp.(*LoginAck)
	assert.True(t, ok)
	assert.Equal(t, "Hello Alvin", loginAck.Welcome)

	// 发送关闭服务请求
	resp, err = client.Request(ACTION_SHUTDOWN, &ShutdownAsk{})
	assert.Nil(t, err)
	assert.Equal(t, &ShutdownAck{}, resp)
}
