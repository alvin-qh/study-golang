package main

import (
	_ "study-gin/app"
	"study-gin/core/conf"
	"study-gin/core/server"
	"study-gin/core/utils"
)

const (
	CONF_FILE = "./application.yaml"
)

func main() {
	// 启动 http 服务
	utils.HttpStart(conf.Config.Server.Address, server.Engine)
}
