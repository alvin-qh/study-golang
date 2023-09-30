package app

import (
	"study-gin/app/routes"
	"study-gin/core/server"
)

// 在此注册其它的路由函数
func init() {
	web := server.Engine.Group("/web")
	{
		web.GET("/user", routes.GetUser)
		web.GET("/user/edit", routes.GetUserEditor)
		web.POST("/user", routes.PostUser)
	}

	api := server.Engine.Group("/api")
	{
		api.GET("/hello", routes.HelloGet)
		api.POST("/hello", routes.HelloPost)
	}
}
