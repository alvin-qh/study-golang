package http

import (
	"context"
	"net/http"
	"time"

	"web/gin/core/utils/signal"

	log "github.com/sirupsen/logrus"
)

// 定义结构体, 从 `http.Server` 结构体继承
type server struct {
	http.Server
}

// 启动 http 服务监听
func (s *server) listen() {
	// 启动监听
	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("cannot start server at \"%v\" caused %v", s.Addr, err)
	}
	log.Infof("server started at \"%v\"", s.Addr)
}

// 关闭 http 服务
func (s *server) shutdown() {
	// 设定一个超时上下文, 超时时间 5 秒
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 关闭 http 服务并等待 5 秒超时
	if err := s.Shutdown(ctx); err != nil {
		log.Fatalf("server was shutdown with error %v", err)
	}
	log.Info("server already shutdown")
}

// 启动 http 服务并等待 `SIGINT` 信号后关闭服务器
//
// 参数:
//   - `address` (`string`): 服务监听地址
//   - `handler` (`http.handler`): `http.handler` 接口对象
func HttpStart(address string, handler http.Handler) {
	// 实例化结构体
	server := &server{
		Server: http.Server{
			Addr:    address, // 监听地址
			Handler: handler, // `http.handler` 接口对象
		},
	}
	log.Infof("starting server at \"%v\"...", server.Addr)

	// 启动协程进行监听
	go server.listen()

	// 在主线程等待进程结束信号
	signal.WaitInterruptSignal()

	log.Infof("shutdown server...")
	// 关闭服务
	server.shutdown()
}
