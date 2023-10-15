package main

import (
	_ "study-golang/gin/app"
	"study-golang/gin/core/conf"
	"study-golang/gin/core/server"
	"study-golang/gin/core/utils/http"
)

const (
	CONF_FILE = "./application.yaml"
)

func main() {
	// 启动 http 服务
	http.HttpStart(conf.Config.Server.Address, server.Engine)
}
