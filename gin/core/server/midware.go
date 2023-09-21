package server

import (
	"strconv"
	"study-gin/core/conf"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func LogMiddleware() gin.HandlerFunc {
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
	headerAllowOrigin      = "Access-Control-Allow-Origin"
	headerAllowMethods     = "Access-Control-Allow-Methods"
	headerAllowHeaders     = "Access-Control-Allow-Headers"
	headerExposeHeaders    = "Access-Control-Expose-Headers"
	headerAllowCredentials = "Access-Control-Allow-Credentials"
	headerMaxAge           = "Access-Control-Max-Age"
)

const headerLogFormat = "\t%v=%v"

// 跨域路由处理
// 参考: https://developer.mozilla.org/zh-CN/docs/Web/HTTP/CORS
func CORSOptionsRoute() gin.HandlerFunc {
	allowMethods := conf.ToString(conf.Config.Server.Cors.AllowMethods, ",")
	log.Infof(headerLogFormat, headerAllowMethods, allowMethods)

	allowHeaders := conf.ToString(conf.Config.Server.Cors.AllowHeaders, ",")
	log.Infof(headerLogFormat, headerAllowHeaders, allowHeaders)

	exposeHeaders := conf.ToString(conf.Config.Server.Cors.ExposeHeaders, ",")
	log.Infof(headerLogFormat, headerExposeHeaders, exposeHeaders)

	allowCredentials := conf.Config.Server.Cors.AllowCredentials
	log.Infof(headerLogFormat, headerAllowCredentials, allowCredentials)

	maxAge := conf.Config.Server.Cors.AllowCredentials
	log.Infof(headerLogFormat, headerMaxAge, maxAge)

	return func(ctx *gin.Context) {
		header := ctx.Writer.Header()

		header.Set(headerAllowMethods, allowMethods)
		header.Set(headerAllowHeaders, allowHeaders)
		header.Set(headerExposeHeaders, exposeHeaders)
		header.Set(headerAllowCredentials, strconv.FormatBool(allowCredentials))
	}
}

// 允许跨域请求的中间件
// 参考: https://developer.mozilla.org/zh-CN/docs/Web/HTTP/CORS
func CORSMiddleware() gin.HandlerFunc {
	log.Info("middleware \"cors\" enabled")
	if !conf.Config.Server.Cors.Enable {
		return func(ctx *gin.Context) { ctx.Next() }
	}

	allowOrigin := conf.ToString(conf.Config.Server.Cors.AllowOrigin, ",")
	log.Infof("\t%v=%v", headerAllowOrigin, allowOrigin)

	return func(ctx *gin.Context) {
		header := ctx.Writer.Header()
		header.Set(headerAllowOrigin, allowOrigin)

		ctx.Next()
	}
}
