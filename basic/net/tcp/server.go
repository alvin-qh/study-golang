package tcp

import (
	"errors"
	"fmt"
	"net"
	"sync/atomic"
	"unsafe"
)

// 定义协议名称
const TCP string = "tcp"

var (
	// 服务器关闭错误
	ErrServerShouldClose = fmt.Errorf("server should close")
)

// 服务端结构体
type Server struct {
	closeCh  chan struct{}    // 服务器关闭事件的 channel
	listener *net.TCPListener // 服务端侦听实例
}

// 启动服务端
func ServerStart(address string) (*Server, error) {
	// 解析服务端监听地址, 形如: "0.0.0.0:8888"
	addr, err := net.ResolveTCPAddr(TCP, address)
	if err != nil {
		return nil, err
	}

	// 监听服务端地址和端口
	listener, err := net.ListenTCP(TCP, addr)
	if err != nil {
		sLog.Fatalf("Network error: %v", err)
		return nil, err
	}
	sLog.Printf("Start listening at %v", addr)

	// 产生服务端对象
	server := &Server{
		closeCh:  make(chan struct{}),
		listener: listener,
	}

	// 调用客户端连接处理函数
	go server.handleAcceptation()

	return server, nil
}

// 接受客户端连接, 启动客户端处理协程
func (s *Server) handleAcceptation() {
	l := s.listener
	if l == nil {
		return
	}

	defer func() {
		close(s.closeCh)
		s.Close()
	}()

	for {
		// 接受一个连接
		conn, err := l.AcceptTCP()
		if err != nil {
			return
		}
		sLog.Printf("New connection coming, %v", conn.RemoteAddr())

		// 处理一次会话
		go s.handleClientSession(conn)
	}
}

// 停止服务端
func (s *Server) Close() {
	l := (*net.TCPListener)(atomic.SwapPointer((*unsafe.Pointer)(unsafe.Pointer(&s.listener)), nil))
	if l != nil {
		l.Close()

		// 等待 accept goroutine 结束
		<-s.closeCh

		sLog.Printf("Server stop successful")
	}
}

// 上下文对象类型
type Context map[string]interface{}

// 服务端请求结构体
type Request struct {
	header  AskHeader // 请求头
	context Context   // 上下文对象
	conn    *TCPConn  // 客户端连接
}

// 处理一次会话
func (s *Server) handleClientSession(conn *net.TCPConn) {
	defer conn.Close()

	// 设置连接保持
	conn.SetKeepAlive(true)

	// 创建请求实例
	req := Request{
		context: make(Context),
		conn:    NewTCPConn(conn),
	}
	defer func() {
		req.conn.Close()
		req.context = nil
	}()

	for {
		// 设置读写超时
		// conn.SetReadDeadline(time.Now().Add(time.Second * 10))
		// conn.SetWriteDeadline(time.Now().Add(time.Second * 30))

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
			if errors.Is(err, ErrServerShouldClose) {
				s.Close()
			} else {
				sLog.Printf("handle session error, %v", err)
			}
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
	sLog.Printf("Do login action, account=%v, password=%v", ask.Account, ask.Password)

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
	sLog.Print("Do shutdown action")

	// 编码响应头
	if err := r.encodeAckHeader(&AckHeader{Action: ACTION_SHUTDOWN, IsOk: true}); err != nil {
		return err
	}

	// 编码响应内容
	if err := r.encodeAckBody(&ShutdownAsk{}); err != nil {
		return err
	}

	return ErrServerShouldClose
}

// 解码请求头
func (r *Request) decodeAskHeader() (AskHeader, error) {
	// 接收请求头数据
	if err := r.conn.Decode(&r.header); err != nil {
		sLog.Printf("Decode ask header failed: %v", err)
		return r.header, err
	}

	sLog.Printf("Ask header received from %v, action=%v", r.conn.RemoteAddr(), r.header.Action)
	return r.header, nil
}

// 解码请求内容
func (r *Request) decodeAskBody(body interface{}) error {
	// 接收请求内容数据
	if err := r.conn.Decode(body); err != nil {
		sLog.Printf("Decode ask body failed: %v", err)
		return err
	}

	sLog.Printf("Ask body received from %v, action=%v", r.conn.RemoteAddr(), r.header.Action)
	return nil
}

// 编码响应头
func (r *Request) encodeAckHeader(header *AckHeader) error {
	// 发送响应头
	if err := r.conn.Encode(header); err != nil {
		sLog.Printf("Encode ack header failed: %v", err)
		return err
	}

	sLog.Printf("Ack header sent to %v, action=%v", r.conn.RemoteAddr(), r.header.Action)
	return nil
}

// 编码相应内容
func (r *Request) encodeAckBody(body interface{}) error {
	// 发送响应内容
	if err := r.conn.Encode(body); err != nil {
		sLog.Printf("Encode ack body failed: %v", err)
	}

	sLog.Printf("Ack body sent to %v, action=%v", r.conn.RemoteAddr(), r.header.Action)
	return nil
}
