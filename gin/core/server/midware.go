package server

import (
	"study-gin/core/conf"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func LogMiddleware(config *conf.Config) gin.HandlerFunc {
	log.Info("middleware \"log\" enabled")
	return func(ctx *gin.Context) {
		startTime := time.Now()
		ctx.Next()
		endTime := time.Now()

		log.Infof("incoming request: \"%v\" visited \"%v %v\", response %v, take %v",
			ctx.ClientIP(),
			ctx.Request.Method,
			ctx.Request.RequestURI,
			ctx.Writer.Status(),
			endTime.Sub(startTime),
		)
	}
}

func JSONMiddleware() gin.HandlerFunc {
	log.Info("middleware \"json\" enabled")
	return func(ctx *gin.Context) {
		header := ctx.Writer.Header()
		header.Set("Content-Type", "application/json; charset=utf-8")
	}
}

const (
	headerAllowOrigin  = "Access-Control-Allow-Origin"
	headerAllowMethods = "Access-Control-Allow-Methods"
)

// 允许跨域请求的中间件
// 参考: https://developer.mozilla.org/zh-CN/docs/Web/HTTP/CORS
func CORSMiddleware(config *conf.Config) gin.HandlerFunc {
	log.Info("middleware \"cors\" enabled")

	allowedOrigins := conf.ToString(config.Server.Cors.AllowedOrigins)
	log.Infof("\t%v=%v", headerAllowOrigin, allowedOrigins)

	allowedMethods := conf.ToString(config.Server.Cors.AllowedMethods)
	log.Infof("\t%v=%v", headerAllowMethods, allowedMethods)

	return func(ctx *gin.Context) {
		header := ctx.Writer.Header()
		header.Set(headerAllowOrigin, conf.ToString(config.Server.Cors.AllowedOrigins))

		if ctx.Request.Method != string(OPTIONS) {
			ctx.Next()
			return
		}
	}
}
