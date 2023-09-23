package server

import (
	"io"

	"github.com/gin-gonic/gin"
)

// 禁用 gin 框架内置的日志
func DisableGinLogger() {
	gin.SetMode(gin.ReleaseMode)

	gin.DefaultWriter = io.Discard
	gin.DefaultErrorWriter = io.Discard
}
