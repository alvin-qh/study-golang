package web

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 测试读取内嵌文件系统的 `index.html` 文件
func TestStartHttpServer(t *testing.T) {
	// 通过设定的 URI 访问文件系统 `asset/index.html` 文件
	// `index.html` 文件为 HTTP 默认文件, 可以省略
	req, err := http.NewRequest("GET", "/asset/", nil)
	assert.Nil(t, err)

	// 创建一个 HTTP 服务路由
	engine := createEngine()

	// 用于测试 HTTP 请求并记录响应
	w := httptest.NewRecorder()

	// 执行测试
	engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

    // 获取响应的 `Body` 内容
	html := w.Body.String()
	assert.Contains(t, html, `<link rel="stylesheet" href="css/main.css">`)
	assert.Contains(t, html, `<title>Go Embed Asset</title>`)
	assert.Contains(t, html, "<h1>Hello, Go Embed</h1>")
	assert.Contains(t, html, `<script type="module" src="js/main.js"></script>`)
}
