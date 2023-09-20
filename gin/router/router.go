package router

import (
	"study-gin/core/server"
	"study-gin/router/routes"

	"github.com/gin-gonic/gin"
)

func Setup(engine *server.Engine) {
	engine.AddRoute(server.GET, "/hello", routes.HelloGet)
	engine.AddRoute(server.POST, "/hello", routes.HelloPost)
	engine.AddRoute(server.OPTIONS, "/*path", func(ctx *gin.Context) {})
}
