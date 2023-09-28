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
	SetupTemplate()
	SetupValidator()

	Engine.Use(
		RecoveryMiddleware(),
		LogMiddleware(),
		CORSMiddleware(),
	)

	if conf.Config.Server.Cors.Enable {
		Engine.OPTIONS("/*path", CORSOptionsRoute())
	}
}
