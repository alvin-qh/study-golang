package main

import (
	"study-gin/app"
	"study-gin/core/conf"
	"study-gin/core/logger"
	"study-gin/core/server"
	"study-gin/core/utils"
)

const (
	CONF_FILE = "./application.yaml"
)

func main() {
	conf.Init(CONF_FILE)
	logger.Init()
	server.Init()
	app.Init()

	utils.StartHttpServer(conf.Config.Server.Address, server.Engine)
}
