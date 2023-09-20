package router

import (
	"study-gin/core/server"
	"study-gin/router/routes"
)

func Setup(engine *server.Engine) {
	engine.AddRoute(server.GET, "/hello", routes.HelloGet)
}
