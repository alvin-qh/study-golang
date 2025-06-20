package main

import (
	_ "study/thirdpart/gin/app"
	"study/thirdpart/gin/core/conf"
	"study/thirdpart/gin/core/server"
	"study/thirdpart/gin/core/utils/http"
)

const (
	CONF_FILE = "./application.yaml"
)

func main() {
	// 启动 http 服务
	http.HttpStart(conf.Config.Server.Address, server.Engine)
}
