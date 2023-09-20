package main

import (
	"study-gin/core/conf"
	"study-gin/core/logger"
	"study-gin/core/server"
	"study-gin/router"
)

const (
	CONF_FILE = "./application.yaml"
)

func main() {
	server.DisableGinLogger()

	conf, err := conf.Load(CONF_FILE)
	if err != nil {
		panic(err)
	}

	logger.Setup(conf)

	engine := server.Create(conf)
	router.Setup(engine)

	engine.Start()
}
