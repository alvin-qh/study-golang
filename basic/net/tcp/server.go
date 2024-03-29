package tcp

import (
	"encoding/gob"
	"fmt"
	"net"
	"sync"
	"time"
)

// 服务端结构体
type Server struct {
	listener *net.TCPListener // 连接监听
	mut      sync.Mutex       // 同步锁
	closeCh  chan struct{}    // 服务器关闭事件的 channel
}

// 获取 Listener 对象
func (s *Server) getListener() *net.TCPListener {
	s.mut.Lock()
	defer s.mut.Unlock()

	return s.listener
}

// 获取服务关闭通知通道对象
func (s *Server) getCloseChan() chan struct{} {
	s.mut.Lock()
	defer s.mut.Unlock()

	return s.closeCh
}

// 启动服务端
func ServerStart(address string) (*Server, error) {
	addr, err := net.ResolveTCPAddr("tcp", address) // 解析服务端监听地址, 形如: "0.0.0.0:8888"
	if err != nil {
		return nil, err
	}

	listener, err := net.ListenTCP("tcp", addr) // 监听服务端地址和端口
	if err != nil {
		lServer.Fatalf("Network error: %v", err)
		return nil, err
	}
	lServer.Printf("Start listening at %v", addr)

	// 产生服务端对象
	server := &Server{
		listener: listener,
		closeCh:  make(chan struct{}),
	}

	// 调用客户端连接处理函数
	go server.handleAcceptation()

	return server, nil
}

// 等待服务端结束
func (s *Server) Join() {
	if ch := s.getCloseChan(); ch != nil {
		<-ch
	}
}

// 停止服务端
func (s *Server) Stop() error {
	var err error

	s.mut.Lock()
	defer s.mut.Unlock()

	// 关闭服务端监听
	if s.listener != nil {
		err = s.listener.Close()
		s.listener = nil
		lServer.Printf("Stop listening")
	}

	// 发送关闭服务端关闭通知
	if s.closeCh != nil {
		close(s.closeCh)
		s.closeCh = nil
		lServer.Printf("Server stop successful")
	}

	return err
}

// 接受客户端连接, 启动客户端处理协程
func (s *Server) handleAcceptation() {
	// 该函数结束后, 表示服务端已结束, 关闭监听并和 channel (发出结束通知)
	defer s.Stop()

	for {
		if listener := s.getListener(); listener != nil {
			// 接受一个连接
			conn, err := s.listener.AcceptTCP()
			if err != nil {
				lServer.Fatalf("Network error: %v", err)
				break
			}
			lServer.Printf("New connection coming, %v", conn.RemoteAddr())

			// 处理一次会话
			go s.handleClientSession(conn)
		} else {
			break
		}
	}
}

// 上下文对象类型
type Context map[string]interface{}

// 创建新的上下文对象
func newContext() Context {
	return make(Context)
}

// 服务端请求结构体
type Request struct {
	server  *Server      // 服务端对象
	conn    *net.TCPConn // 客户端连接
	header  AskHeader    // 请求头
	context Context      // 上下文对象
}

// 解码请求头
func (r *Request) decodeAskHeader() (AskHeader, error) {
	decoder := gob.NewDecoder(r.conn)

	// 接收请求头数据
	if err := decoder.Decode(&r.header); err != nil {
		lServer.Printf("Decode ask header failed: %v", err)
		return r.header, err
	}

	lServer.Printf("Ask header received from %v, action=%v", r.conn.RemoteAddr(), r.header.Action)
	return r.header, nil
}

// 解码请求内容
func (r *Request) decodeAskBody(body interface{}) error {
	decoder := gob.NewDecoder(r.conn)

	// 接收请求内容数据
	if err := decoder.Decode(body); err != nil {
		lServer.Printf("Decode ask body failed: %v", err)
		return err
	}

	lServer.Printf("Ask body received from %v, action=%v", r.conn.RemoteAddr(), r.header.Action)
	return nil
}

// 编码响应头
func (r *Request) encodeAckHeader(header *AckHeader) error {
	encoder := gob.NewEncoder(r.conn)

	// 发送响应头
	if err := encoder.Encode(header); err != nil {
		lServer.Printf("Encode ack header failed: %v", err)
		return err
	}

	lServer.Printf("Ack header sent to %v, action=%v", r.conn.RemoteAddr(), r.header.Action)
	return nil
}

// 编码相应内容
func (r *Request) encodeAckBody(body interface{}) error {
	encoder := gob.NewEncoder(r.conn)

	// 发送响应内容
	if err := encoder.Encode(body); err != nil {
		lServer.Printf("Encode ack body failed: %v", err)
	}

	lServer.Printf("Ack body sent to %v, action=%v", r.conn.RemoteAddr(), r.header.Action)
	return nil
}

// cspell: ignore sess
// 处理一次会话
func (s *Server) handleClientSession(conn *net.TCPConn) {
	defer conn.Close()

	// 设置连接保持
	conn.SetKeepAlive(true)
	for {
		// 设置读写超时
		conn.SetReadDeadline(time.Now().Add(time.Second * 10))
		conn.SetWriteDeadline(time.Now().Add(time.Second * 30))

		req := Request{
			server:  s,
			conn:    conn,
			context: newContext(),
		}

		// 解码请求头
		header, err := req.decodeAskHeader()
		if err != nil {
			break
		}

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
	// 解码请求内容
	ask := LoginAsk{}
	if err := r.decodeAskBody(&ask); err != nil {
		return err
	}
	lServer.Printf("Do login action, account=%v, password=%v", ask.Account, ask.Password)

	// 编码响应头
	if err := r.encodeAckHeader(&AckHeader{Action: ACTION_LOGIN, IsOk: true}); err != nil {
		return err
	}

	// 编码响应内容
	content := fmt.Sprintf("Hello %v", ask.Account)
	if err := r.encodeAckBody(&LoginAck{Welcome: content}); err != nil {
		return err
	}

	return nil
}

// 处理关闭服务会话
func (r *Request) handleShutdownAction() error {
	// 解码请求内容
	ask := ShutdownAsk{}
	if err := r.decodeAskBody(&ask); err != nil {
		return err
	}
	lServer.Print("Do shutdown action")

	// 编码响应头
	if err := r.encodeAckHeader(&AckHeader{Action: ACTION_SHUTDOWN, IsOk: true}); err != nil {
		return err
	}

	// 编码响应内容
	if err := r.encodeAckBody(&ShutdownAsk{}); err != nil {
		return err
	}

	// 关闭服务器
	r.server.Stop()
	return nil
}
