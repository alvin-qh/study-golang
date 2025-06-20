package server

import (
	"errors"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"strconv"
	"strings"
	"time"

	"web/gin/core/conf"
	"web/gin/core/utils/callstack"
	"web/gin/core/utils/value"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// 定义日志记录中间件函数
//
// 返回:
//   - gin 框架中间件函数
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

		for _, err := range ctx.Errors {
			log.Errorf("error caused when visit \"%v\", error: %v", ctx.Request.URL.Path, err.Err)
		}
	}
}

// 定义异常恢复中间件
//
// 当代码中抛出 `panic` 错误时, 该中间件负责捕获并输出日志
//
// 输出的日志包括发生 `panic` 错误的原因和调用堆栈
//
// 参数:
//   - `handlers` (`...gin.RecoveryFunc`): 其它处理异常恢复的回调
//
// 返回
//   - 中间件函数
func RecoveryMiddleware(handles ...gin.RecoveryFunc) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 判断连接是否中断
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					var se *os.SyscallError
					if errors.As(ne, &se) {
						seStr := strings.ToLower(se.Error())
						if strings.Contains(seStr, "broken pipe") ||
							strings.Contains(seStr, "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				if brokenPipe {
					// 转存请求对象
					req, _ := httputil.DumpRequest(ctx.Request, false)
					// 获取所有请求头
					headers := strings.Split(string(req), "\r\n")

					// 遍历请求头, 过滤掉 `Authorization` 头的内容
					for idx, header := range headers {
						if strings.HasPrefix(header, "Authorization") {
							current := strings.Split(header, ":")
							headers[idx] = current[0] + ": *"
						}
					}

					// 将请求头连接为一个字符串
					headersToStr := strings.Join(headers, "\r\n")

					// 输出 panic 错误以及请求头内容
					log.Errorf("%s\n%s", err, headersToStr)
				} else {
					// 输出 panic 错误以及调用栈信息
					log.Errorf("[Recovery] panic recovered: %s\n%s", err, callstack.CallStack(3))
				}

				if brokenPipe {
					// 如果连接中断, 则输出错误信息并中断请求, 无法输出响应状态
					ctx.Error(err.(error))
					ctx.Abort()
				} else {
					// 如果连接未中断, 则输出响应状态
					ctx.AbortWithStatus(http.StatusInternalServerError)

					// 调用后续的处理函数 (如果存在)
					for _, h := range handles {
						h(ctx, err)
					}
				}
			}
		}()

		ctx.Next()
	}
}

// 定义跨域的 http 请求头常量
const (
	headerOrigin           = "Origin"
	headerAllowOrigin      = "Access-Control-Allow-Origin"
	headerAllowMethods     = "Access-Control-Allow-Methods"
	headerAllowHeaders     = "Access-Control-Allow-Headers"
	headerExposeHeaders    = "Access-Control-Expose-Headers"
	headerAllowCredentials = "Access-Control-Allow-Credentials"
	headerMaxAge           = "Access-Control-Max-Age"
)

// 输出 http 请求头的日志模板
const headerLogFormat = "\t%v=%v"

// 定义跨域处理路由函数
//
// 关于跨域, 参考: https://developer.mozilla.org/zh-CN/docs/Web/HTTP/CORS
//
// 返回:
//   - gin 框架路由处理函数
func CORSOptionsRoute() gin.HandlerFunc {
	allowMethods := value.JoinAny(conf.Config.Server.Cors.AllowMethods, ",")
	log.Infof(headerLogFormat, headerAllowMethods, allowMethods)

	allowHeaders := value.JoinAny(conf.Config.Server.Cors.AllowHeaders, ",")
	log.Infof(headerLogFormat, headerAllowHeaders, allowHeaders)

	exposeHeaders := value.JoinAny(conf.Config.Server.Cors.ExposeHeaders, ",")
	log.Infof(headerLogFormat, headerExposeHeaders, exposeHeaders)

	maxAge := conf.Config.Server.Cors.MaxAge
	log.Infof(headerLogFormat, headerMaxAge, maxAge)

	return func(ctx *gin.Context) {
		header := ctx.Writer.Header()
		header.Set(headerAllowMethods, allowMethods)
		header.Set(headerAllowHeaders, allowHeaders)
		header.Set(headerExposeHeaders, exposeHeaders)
		header.Set(headerMaxAge, strconv.FormatInt(int64(maxAge), 10))
	}
}

// 定义跨域处理中间件函数
//
// 关于跨域, 参考: https://developer.mozilla.org/zh-CN/docs/Web/HTTP/CORS
//
// 返回:
//   - gin 框架中间件函数
func CORSMiddleware() gin.HandlerFunc {
	// 如果未启用跨域, 则返回一个空中间件函数, 不做任何处理
	if !conf.Config.Server.Cors.Enable {
		return func(ctx *gin.Context) { ctx.Next() }
	}

	log.Info("middleware \"cors\" enabled")

	// 获取配置的跨域 HTTP 头设置
	allowOrigin := MakeAllowOrigin(conf.Config.Server.Cors.AllowOrigin)
	log.Infof("\t%v=%v", headerAllowOrigin, allowOrigin)

	allowCredentials := strconv.FormatBool(conf.Config.Server.Cors.AllowCredentials)
	log.Infof(headerLogFormat, headerAllowCredentials, allowCredentials)

	// 返回中间件函数, 用于为所有响应增加跨域 HTTP 头
	return func(ctx *gin.Context) {
		origin := ctx.Request.Header.Get(headerOrigin)

		if len(origin) > 0 {
			header := ctx.Writer.Header()

			_, ok := allowOrigin["*"]
			if !ok {
				_, ok = allowOrigin[origin]
			}

			if ok {
				header.Set(headerAllowOrigin, origin)
				header.Set(headerAllowCredentials, allowCredentials)
			} else {
				ctx.AbortWithStatus(http.StatusForbidden)
				log.Infof("invalid origin \"%v\", forbidden", origin)
			}
		}
		ctx.Next()
	}
}
