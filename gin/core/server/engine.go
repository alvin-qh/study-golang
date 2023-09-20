package server

import (
	"study-gin/core/conf"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type Engine struct {
	engine  *gin.Engine
	address []string
}

type Method string

// 定义请求方法常量
const (
	GET     Method = "GET"
	PUT     Method = "PUT"
	POST    Method = "POST"
	DELETE  Method = "DELETE"
	HEAD    Method = "HEAD"
	PATCH   Method = "PATCH"
	OPTIONS Method = "OPTIONS"
)

func Create(config *conf.Config) *Engine {
	engine := gin.New()
	engine.Use(gin.Recovery(), LogMiddleware(config), JSONMiddleware(), CORSMiddleware(config))

	return &Engine{
		engine:  engine,
		address: config.Server.Address,
	}
}

func (e *Engine) AddRoute(method Method, path string, handlers ...gin.HandlerFunc) {
	e.engine.Handle(string(method), path, handlers...)
	log.Infof("route was add, method \"%v\", path \"%v\"", method, path)
}

func (e *Engine) Start() {
	log.Infof("server start at \"%v\"", e.address)

	err := e.engine.Run(e.address...)
	if err != nil {
		log.Errorf("cannot start server, error is \"%v\"", err)
	}
}
