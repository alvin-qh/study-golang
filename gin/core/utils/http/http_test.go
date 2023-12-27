package http

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

const (
	TEST_PORT = "15000"
)

// 测试 `HttpStart` 函数
//
// 验证 http 服务的启动和停止
func TestStartAndShutdown(t *testing.T) {
	// 定义发送协程结束消息的通道
	ch := make(chan struct{}, 1)

	// 启动协程, 在协程中启动 http 服务器, 并等待 `SIGINT` 信号后停止服务器
	go func() {
		// 实例化 gin 引擎对象
		engine := gin.New()

		// 添加一个 GET 处理函数
		engine.GET("/", func(ctx *gin.Context) {
			// 返回一个字符串
			ctx.String(200, "Hello")
		})

		// 启动 http 服务器, 并等待 `SIGINT` 信号后停止服务器
		HttpStart(fmt.Sprintf(":%v", TEST_PORT), engine)

		// http 服务器停止后, 发送协程结束事件
		ch <- struct{}{}
	}()

	time.Sleep(500 * time.Millisecond)

	// 访问测试地址
	resp, err := http.Get(fmt.Sprintf("http://localhost:%v", TEST_PORT))
	assert.NoError(t, err)

	// 读取响应结果
	r := bufio.NewReader(resp.Body)
	s, err := r.ReadString('\n')
	assert.ErrorIs(t, err, io.EOF)

	// 确认响应结果
	assert.Equal(t, "Hello", s)

	// 获取当前进程对象
	process, err := os.FindProcess(syscall.Getpid())
	assert.NoError(t, err)

	// 发送停止信号
	err = process.Signal(os.Interrupt)
	assert.NoError(t, err)

	<-ch
}
