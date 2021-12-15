package tcp

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func TestTcpConnect(t *testing.T) {
	// 启动服务器
	server, err := StartTCPServer("0.0.0.0:8888")
	assert.NoError(t, err)
	defer server.StopTCPServer()

	// 连接服务器
	client, err := TCPConnect("127.0.0.1:8888")
	assert.NoError(t, err)
	defer client.Close()

	// 发送登录请求
	resp, err := client.Request(ACTION_LOGIN, &TCPLoginAsk{
		Account:  "Alvin",
		Password: "password",
	})
	assert.NoError(t, err)

	// 接收登录响应
	loginAck, ok := resp.(*TCPLoginAck)
	assert.True(t, ok)
	assert.Equal(t, "Hello Alvin", loginAck.Welcome)

	// 发送关闭服务请求
	resp, err = client.Request(ACTION_SHUTDOWN, &TCPShutdownAsk{})
	assert.NoError(t, err)
	assert.Equal(t, &TCPShutdownAck{}, resp)

	server.Join()
}
