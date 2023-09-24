package app

import (
	"study-gin/app/routes"
	"study-gin/core/server"
)

// 在此注册其它的路由函数
func Init() {
	web := server.Engine.Group("/web")
	{
		web.GET("/render", routes.RenderHTML)
	}

	api := server.Engine.Group("/api")
	{
		api.GET("/hello", routes.HelloGet)
		api.POST("/hello", routes.HelloPost)
	}
}
