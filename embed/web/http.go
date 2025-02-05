package web

import (
	"embed"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

var (
	//go:embed asset/*
	STATIC_FS embed.FS
)

const (
	staticPathPrefix   = "/asset/"
	cacheControlHeader = "Cache-Control"
	cacheControlValue  = "public, max-age=2592000" // 缓存一个月
	eTagHeader         = "ETag"
	ifNoneMatchHeader  = "If-None-Match"
)

func createEngine() *gin.Engine {
	engine := gin.Default()

	engine.Use(func(c *gin.Context) {
		// 如果请求的 URI 以 `/asset/` 开头, 则设置缓存头
		if strings.HasPrefix(c.Request.URL.Path, "/asset/") {
			c.Header(cacheControlHeader, cacheControlValue)

			fileName := c.Request.URL.Path[1:]
			f, err := STATIC_FS.Open(fileName)
			if err == nil {
				st, err := f.Stat()
				if err != nil {
					c.Abort()
					return
				}
				eTag := fmt.Sprintf("%v-%v", base64.StdEncoding.EncodeToString([]byte(st.Name())), st.Size())
				c.Header(eTagHeader, eTag)

				if match := c.GetHeader(ifNoneMatchHeader); match == eTag {
					if strings.Contains(match, eTag) {
						c.Status(http.StatusNotModified)
						c.Abort()
						return
					}
				}
			}
		}
		c.Next()
	})

	// `gin` 框架官方支持内嵌文件系统的方法
	// 缺点是: 文件系统的路径会成为 URL 的一部分, 例如: 要通过 `/asset/asset/index.html` 来访问文件
	// (如果将 URI 参数设置为 `/`, 可以通过 `/asset/index.html` 访问文件, 但同时会占用 `/` 路由, 不推荐)
	// engine.StaticFS("/asset", http.FS(STATIC_FS))

	// 映射文件系统服务
	staticFileSrv := http.FileServer(http.FS(STATIC_FS))

	// 可以通过如下方式解决 URI 和文件路径重复的问题, 可以通过 `/asset/index.html` 访问到文件,
	// 这里 URI 中的 `/asset` 必须和文件系统的根目录同名, 即文件系统必须映射为 `//go:embed asset`
	engine.Any("/asset/*filepath", func(c *gin.Context) {
		// 将请求和相应定向到服务
		staticFileSrv.ServeHTTP(c.Writer, c.Request)
	})

	return engine
}

// 启动 HTTP 服务
func StartServer(addr string) {
	engine := createEngine()
	engine.Run(addr)
}
