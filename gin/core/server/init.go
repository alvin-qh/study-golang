package server

import (
	"study-gin/core/conf"

	"github.com/gin-gonic/gin"
)

var (
	Engine = gin.New()
)

func Init() {
	DisableGinLogger()

	Engine.Use(gin.Recovery(), LogMiddleware(), JSONMiddleware(), CORSMiddleware())
	if conf.Config.Server.Cors.Enable {
		Engine.OPTIONS("/*path", CORSOptionsRoute())
	}
}
