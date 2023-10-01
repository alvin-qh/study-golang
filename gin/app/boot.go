package app

import (
	"study-gin/app/routes"
	"study-gin/core/server"
)

// 在此注册其它的路由函数
func init() {
	web := server.Engine.Group("/web")
	{
		web.GET("/user", routes.WebGetUser)
		web.GET("/user/edit", routes.WebGetUserEditor)
		web.POST("/user", routes.WebPostUser)
	}

	api := server.Engine.Group("/api")
	{
		api.GET("/user", routes.ApiGetUser)
		api.POST("/user", routes.ApiPostUser)
	}
}
