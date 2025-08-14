package main

import (
	_ "study/web/gin/app"
	"study/web/gin/core/conf"
	"study/web/gin/core/server"
	"study/web/gin/core/utils/http"
)

const (
	CONF_FILE = "./application.yaml"
)

func main() {
	// 启动 http 服务
	http.HttpStart(conf.Config.Server.Address, server.Engine)
}
