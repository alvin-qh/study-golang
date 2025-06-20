package main

import (
	_ "web/gin/app"
	"web/gin/core/conf"
	"web/gin/core/server"
	"web/gin/core/utils/http"
)

const (
	CONF_FILE = "./application.yaml"
)

func main() {
	// 启动 http 服务
	http.HttpStart(conf.Config.Server.Address, server.Engine)
}
