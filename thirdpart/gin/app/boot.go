package app

import (
	"net/http"
	"study/thirdpart/gin/app/routes"
	"study/thirdpart/gin/core/conf"
	"study/thirdpart/gin/core/server"
	"study/thirdpart/gin/core/utils/maptostruct"

	"github.com/gin-gonic/gin"
)

var (
	MapToStruct    = maptostruct.New("json")
	allHttpMethods = []string{
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodHead,
		http.MethodDelete,
	}
)

func AllMethod(group *gin.RouterGroup, relativePath string, handlers ...gin.HandlerFunc) {
	for _, m := range allHttpMethods {
		group.Handle(m, relativePath, handlers...)
	}
}

// 在此注册其它的路由函数
func init() {
	web := server.Engine.Group("/web")
	{
		web.GET("/users", routes.WebGetUsers)
		web.GET("/users/new", routes.WebNewUsers)
		web.POST("/users", routes.WebPostUsers)
		web.GET("/users/:id", routes.WebEditUsers)
		web.POST("/users/:id", routes.WebPutUsers)
	}

	api := server.Engine.Group("/api")
	{
		api.GET("/users", routes.ApiGetUsers)
		api.POST("/users", routes.ApiPostUsers)
		api.GET("/users/:id", routes.ApiGetUserById)
	}

	proxy := server.Engine.Group("/proxy")
	{
		AllMethod(proxy, "/*path", routes.Proxy)
	}
}

func init() {
	// 如果启用跨域配置, 则对 OPTIONS 请求启用响应
	if conf.Config.Server.Cors.Enable {
		// 对所有 OPTIONS 请求进行处理
		server.Engine.OPTIONS("/*path", server.CORSOptionsRoute())
	}
}
