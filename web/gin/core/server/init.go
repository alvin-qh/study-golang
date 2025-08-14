package server

import (
	_ "study/web/gin/core/logger"

	"github.com/gin-gonic/gin"
)

var (
	Engine = gin.New()
)

// 初始化 http 服务
func init() {
	// 关闭 gin 本身的日志
	DisableGinLogger()
	// 设置 gin HTML 模板
	SetupTemplate()
	// 设置 gin 验证引擎
	SetupValidator()

	// 设置 gin 中间件
	Engine.Use(
		RecoveryMiddleware(), // recover 中间件, 用于从异常中恢复
		LogMiddleware(),      // 记录日志的中间件
		CORSMiddleware(),     // 跨域中间件
	)
}
