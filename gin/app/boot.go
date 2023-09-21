package app

import (
	"study-gin/app/routes"
	"study-gin/core/server"
)

// 在此注册其它的路由函数
func Init() {
	server.Engine.GET("/hello", routes.HelloGet)
	server.Engine.POST("/hello", routes.HelloPost)
}
