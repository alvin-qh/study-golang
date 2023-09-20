package server

import (
	"strconv"
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
	headerAllowOrigin      = "Access-Control-Allow-Origin"
	headerAllowMethods     = "Access-Control-Allow-Methods"
	headerAllowHeaders     = "Access-Control-Allow-Headers"
	headerExposeHeaders    = "Access-Control-Expose-Headers"
	headerAllowCredentials = "Access-Control-Allow-Credentials"
	headerMaxAge           = "Access-Control-Max-Age"
)

// 允许跨域请求的中间件
// 参考: https://developer.mozilla.org/zh-CN/docs/Web/HTTP/CORS
func CORSMiddleware(config *conf.Config) gin.HandlerFunc {
	enable := config.Server.Cors.Enable
	if !enable {
		return func(ctx *gin.Context) {
			ctx.Next()
		}
	} else {
		log.Info("middleware \"cors\" enabled")

		allowOrigin := conf.ToString(config.Server.Cors.AllowOrigin, ",")
		log.Infof("\t%v=%v", headerAllowOrigin, allowOrigin)

		allowMethods := conf.ToString(config.Server.Cors.AllowMethods, ",")
		log.Infof("\t%v=%v", headerAllowMethods, allowMethods)

		allowHeaders := conf.ToString(config.Server.Cors.AllowHeaders, ",")
		log.Infof("\t%v=%v", headerAllowHeaders, allowHeaders)

		exposeHeaders := conf.ToString(config.Server.Cors.ExposeHeaders, ",")
		log.Infof("\t%v=%v", headerExposeHeaders, exposeHeaders)

		allowCredentials := config.Server.Cors.AllowCredentials
		log.Infof("\t%v=%v", headerAllowCredentials, allowCredentials)

		maxAge := config.Server.Cors.AllowCredentials
		log.Infof("\t%v=%v", headerMaxAge, maxAge)

		return func(ctx *gin.Context) {
			header := ctx.Writer.Header()
			header.Set(headerAllowOrigin, allowOrigin)

			if ctx.Request.Method == string(OPTIONS) {
				header.Set(headerAllowMethods, allowMethods)
				header.Set(headerAllowHeaders, allowHeaders)
				header.Set(headerExposeHeaders, exposeHeaders)
				header.Set(headerAllowCredentials, strconv.FormatBool(allowCredentials))

				ctx.AbortWithStatus(200)
				return
			}

			ctx.Next()
		}
	}
}
