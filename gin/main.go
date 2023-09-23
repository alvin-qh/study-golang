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
	// 初始化配置信息
	conf.Init(CONF_FILE)

	// 初始化日志
	logger.Init()

	// 初始化 http 服务
	server.Init()

	// 初始化应用程序
	app.Init()

	// 启动 http 服务
	utils.HttpStart(conf.Config.Server.Address, server.Engine)
}
