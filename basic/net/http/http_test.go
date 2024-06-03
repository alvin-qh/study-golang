package http

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// 测试 HTTP 请求
func TestHttp_Network(t *testing.T) {
	// 实例化 Http 服务端路由
	mux := http.NewServeMux()

	// 实例化 Http 服务端
	srv := &http.Server{
		Addr:    ":18888", // 端口号
		Handler: mux,      // 路由实例
	}
	defer srv.Shutdown(context.Background())

	// 通知服务端关闭的信道
	cShutdown := make(chan struct{})

	// 启动服务端的 goroutine
	go func() {
		// 启动服务器监听
		err := srv.ListenAndServe()
		// 服务端结束, 确认返回错误信息
		assert.EqualError(t, err, http.ErrServerClosed.Error())

		// 关闭信道, 通知服务端已关闭
		close(cShutdown)
	}()

	// 添加服务端路由和处理函数
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// 获取 URL 上的查询参数
		qs := r.URL.Query()
		name := qs.Get("name")

		// 添加响应 Http 头信息
		h := w.Header()
		h.Add("X-User", name)

		// 写入 Http 响应内容
		w.Write([]byte(fmt.Sprintf("Hello, %v!", name)))
	})

	success := false

	// 启动客户端 goroutine
	go func() {
		// 生成访问 URL
		uri := fmt.Sprintf("http://127.0.0.1:18888?name=%v", url.QueryEscape("Alvin"))

		// 最多重试 10 次, 访问服务端
		for i := 0; i < 10 && !success; i++ {
			// 通过指定 URL 向服务端发送请求, 获取响应结果
			resp, err := http.Get(uri)
			if err != nil {
				// 如果服务端尚未响应, 则等待 10ms 后重试
				time.Sleep(10 * time.Millisecond)
				continue
			}

			// 获取响应结果的 HTTP 头信息, 确认指定信息存在
			header := resp.Header.Get("X-User")
			assert.Equal(t, "Alvin", header)

			// 获取响应内容, 确认响应内容符合预期
			body, err := io.ReadAll(resp.Body)
			assert.Nil(t, err)
			assert.Equal(t, "Hello, Alvin!", string(body))

			success = true
		}

		// 完成测试, 关闭服务端
		srv.Shutdown(context.Background())
	}()

	// 等待服务端关闭
	<-cShutdown

	assert.True(t, success)
}
