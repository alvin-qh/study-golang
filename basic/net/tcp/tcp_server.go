package tcp

import (
	"encoding/gob"
	"fmt"
	"net"
	"sync"
	"time"
)

// 服务端结构体
type TCPServer struct {
	listener *net.TCPListener // 连接监听
	mut      sync.Mutex       // 同步锁
	closeCh  chan struct{}    // 服务器关闭事件的 channel
}

// 获取 Listener 对象
func (s *TCPServer) getListener() *net.TCPListener {
	s.mut.Lock()
	defer s.mut.Unlock()

	return s.listener
}

// 获取服务关闭通知通道对象
func (s *TCPServer) getCloseChan() chan struct{} {
	s.mut.Lock()
	defer s.mut.Unlock()

	return s.closeCh
}

// 启动服务端
func StartTCPServer(address string) (*TCPServer, error) {
	addr, err := net.ResolveTCPAddr("tcp", address) // 解析服务端监听地址，形如："0.0.0.0:8888"
	if err != nil {
		return nil, err
	}

	listener, err := net.ListenTCP("tcp", addr) // 监听服务端地址和端口
	if err != nil {
		logger.Fatalf("network error: %v", err)
		return nil, err
	}

	// 产生服务端对象
	server := &TCPServer{
		listener: listener,
		closeCh:  make(chan struct{}),
	}

	// 调用客户端连接处理函数
	go server.handleAcception()

	return server, nil
}

// 等待服务端结束
func (s *TCPServer) Join() {
	if ch := s.getCloseChan(); ch != nil {
		<-ch
	}
}

// 停止服务端
func (s *TCPServer) StopTCPServer() error {
	var err error

	s.mut.Lock()
	defer s.mut.Unlock()

	// 关闭服务端监听
	if s.listener != nil {
		err = s.listener.Close()
		s.listener = nil
	}

	// 发送关闭服务端关闭通知
	if s.closeCh != nil {
		close(s.closeCh)
		s.closeCh = nil
	}

	return err
}

// 接受客户端连接，启动客户端处理协程
func (s *TCPServer) handleAcception() {
	if s.listener == nil {
		logger.Fatalf("Server started already")
		return
	}

	// 该函数结束后，表示服务端已结束，关闭监听并和 channel（发出结束通知）
	defer s.StopTCPServer()

	if listener := s.getListener(); listener != nil {
		logger.Printf("Ready to accepting, %v", s.listener.Addr())
	}

	for {
		if listener := s.getListener(); listener != nil {
			// 接受一个连接
			conn, err := s.listener.AcceptTCP()
			if err != nil {
				logger.Fatalf("Network error: %v", err)
			}
			// 处理一次会话
			go s.handleClientSession(conn)
		} else {
			break
		}
	}
}

// 服务端请求结构体
type Request struct {
	server  *TCPServer   // 服务端对象
	conn    *net.TCPConn // 客户端连接
	decoder *gob.Decoder // 对象编码解码对象
	encoder *gob.Encoder
	header  TCPAskHeader           // 请求头
	context map[string]interface{} // 上下文对象
}

// 解码请求头
func (r *Request) decodeAskHeader() (TCPAskHeader, error) {
	if err := r.decoder.Decode(&r.header); err != nil {
		logger.Printf("Decode ask header failed: %v", err)
		return r.header, err
	}
	return r.header, nil
}

// 解码请求内容
func (r *Request) decodeAskBody(body interface{}) error {
	if err := r.decoder.Decode(body); err != nil {
		logger.Printf("Decode ask body failed: %v", err)
		return err
	}
	return nil
}

// 编码响应头
func (r *Request) encodeAckHeader(header *TCPAckHeader) error {
	if err := r.encoder.Encode(header); err != nil {
		logger.Printf("Encode ack header failed: %v", err)
		return err
	}
	return nil
}

// 编码相应内容
func (r *Request) encodeAckBody(body interface{}) error {
	if err := r.encoder.Encode(body); err != nil {
		logger.Printf("Encode ack body failed: %v", err)
	}
	return nil
}

// cspell: ignore sess
// 处理一次会话
func (s *TCPServer) handleClientSession(conn *net.TCPConn) {
	logger.Printf("New connection comming, %v", conn.RemoteAddr())
	defer conn.Close()

	// 设置连接保持
	conn.SetKeepAlive(true)
	for {
		// 设置读写超时
		conn.SetReadDeadline(time.Now().Add(time.Second * 10))
		conn.SetWriteDeadline(time.Now().Add(time.Second * 30))

		req := &Request{
			server:  s,
			conn:    conn,
			decoder: gob.NewDecoder(conn),
			encoder: gob.NewEncoder(conn),
			context: make(map[string]interface{}),
		}

		header, err := req.decodeAskHeader()
		if err != nil {
			break
		}

		logger.Printf("New package header received, action=%v", header.Action)

		// 根据请求头处理会话
		switch header.Action {
		case ACTION_LOGIN:
			err = req.handleLoginAction()
		case ACTION_SHUTDOWN:
			err = req.handleShutdownAction()
		}

		if err != nil {
			break
		}
	}
}

// 处理登录会话
func (r *Request) handleLoginAction() error {
	ask := TCPLoginAsk{}
	if err := r.decodeAskBody(&ask); err != nil {
		return err
	}
	logger.Printf("New ACTION_LOGIN body received, account=%v, password=%v", ask.Account, ask.Password)

	if err := r.encodeAckHeader(&TCPAckHeader{
		Action: ACTION_LOGIN,
		IsOk:   true,
	}); err != nil {
		return err
	}

	if err := r.encodeAckBody(&TCPLoginAck{
		Welcome: fmt.Sprintf("Hello %v", ask.Account),
	}); err != nil {
		return err
	}

	return nil
}

// 处理关闭服务会话
func (r *Request) handleShutdownAction() error {
	ask := TCPShutdownAsk{}
	if err := r.decodeAskBody(&ask); err != nil {
		return err
	}
	logger.Printf("New ACTION_SHUTDOWN body received")

	if err := r.encodeAckHeader(&TCPAckHeader{
		Action: ACTION_SHUTDOWN,
		IsOk:   true,
	}); err != nil {
		return err
	}

	if err := r.encodeAckBody(&TCPShutdownAsk{}); err != nil {
		return err
	}

	r.server.StopTCPServer()
	return nil
}
