package app

import (
	"study-gin/app/routes"
	"study-gin/core/server"
	"study-gin/core/utils/maptostruct"
)

var (
	MapToStruct = maptostruct.New("json")
)

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
}
