package udp

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/google/uuid"
)

// 定义错误值
var (
	ErrServerClosed     = errors.New("server closed")
	ErrInvalidPackage   = errors.New("invalid package")
	ErrInvalidSessionId = errors.New("invalid session id")
	ErrCloseServer      = errors.New("client ask close server")
)

// 定义 Server 结构体
type Server struct {
	conn    *net.UDPConn  // UDP 连接对象
	sendCh  chan Response // 发送相应包的 channel
	session sync.Map      // 保存 session 信息
	closeCh chan struct{} // 发送服务关闭消息的 channel
}

// 保存 Session 信息
func (s *Server) saveSession(sessionId SessionId) {
	s.session.Store(sessionId, struct{}{})
}

// 判断 Session 是否存在
func (s *Server) hasSession(sessionId SessionId) bool {
	_, ok := s.session.Load(sessionId)
	return ok
}

// 启动服务器
func ServerStart(address string) (*Server, error) {
	// 解析监听地址
	addr, err := net.ResolveUDPAddr("udp4", address)
	if err != nil {
		lServer.Fatalf("Cannot resolve address %v", address)
		return nil, err
	}

	// 监听指定端口和地址
	conn, err := net.ListenUDP("udp4", addr)
	if err != nil {
		lServer.Fatalf("Cannot Listen UDP at %v", address)
		return nil, err
	}
	lServer.Printf("Server was listened at %v", addr)

	// 创建服务端对象
	server := &Server{
		conn:    conn,
		sendCh:  make(chan Response, 100),
		session: sync.Map{},
		closeCh: make(chan struct{}),
	}

	// 接收信息
	go server.handleReceiveMessage()

	// 发送信息
	go server.handleSendMessage()

	defer func() {
		r := recover()
		if err, ok := r.(error); ok {
			lServer.Panicf("Panic %v", err)
		}
	}()

	return server, nil
}

// 停止服务端
func (s *Server) Stop() error {
	var err error

	// 关闭发送通道
	if s.sendCh != nil {
		close(s.sendCh)
		s.sendCh = nil
	}

	// 关闭连接
	if s.conn != nil {
		err = s.conn.Close()
		s.conn = nil
	}

	return err
}

// 等待服务器关闭
func (s *Server) Join() {
	ch := s.closeCh
	if ch != nil {
		<-ch
	}
}

// 处理发送信息
func (s *Server) handleSendMessage() {
	// 获取发送 channel 对象
	ch := s.sendCh
	if ch == nil {
		return
	}

	// 从 channel 中获取待发送的数据
	for resp := range ch {
		// 获取 UDP 连接对象
		conn := s.conn
		if conn == nil {
			break
		}

		// 将数据写入 UDP 连接
		buf := bytes.NewBuffer(make([]byte, 0, PACKAGE_LIMIT))

		encoder := gob.NewEncoder(buf)

		// 将发送数据编码后写入缓冲
		if err := encoder.Encode(resp.pack); err != nil {
			lServer.Fatalf("Cannot send response, caused %v", err)
		}

		// 缓冲数据写入 UDP
		if n, err := conn.WriteToUDP(buf.Bytes(), resp.addr); err != nil {
			lServer.Fatalf("Cannot send response, caused %v", err)
			break
		} else {
			lServer.Printf("%v bytes write to %v", n, resp.addr)
		}

		// 发送完数据包后, 如果需要关闭服务器, 则关闭服务器, 防止最后一个数据包无法发送
		if resp.close {
			s.Stop()
		}
	}
}

// 处理接收数据
func (s *Server) handleReceiveMessage() {
	ch := s.sendCh
	if ch == nil {
		return
	}

	for {
		// 接收数据
		req, err := s.receivePackage()
		if err != nil {
			lServer.Printf("Cannot receive package, caused: %v", err)
			break
		}
		lServer.Printf("Received package from %v, action=%v", req.addr, req.action)

		// 处理已接收的数据
		go func() {
			var pack Package

			closeServer := false

			switch req.action {
			case ACTION_LOGIN:
				// 处理登录消息
				pack, err = req.handleActionLogin()
			case ACTION_SHUTDOWN:
				// 处理关闭服务消息
				pack, err = req.handleActionShutdown()
			}

			if err != nil {
				if errors.Is(err, ErrCloseServer) {
					closeServer = true
				} else {
					log.Fatal(err)
				}
			}

			// 保存 session 对象
			s.saveSession(req.sessionId)

			// 发送数据到 channel
			ch <- Response{
				addr:  req.addr,
				pack:  pack,
				close: closeServer,
			}
		}()
	}

	// 发送关闭消息
	close(s.closeCh)
	s.closeCh = nil
}

// 接收数据包
func (s *Server) receivePackage() (*Request, error) {
	// 获取 UDP 连接
	conn := s.conn
	if conn == nil {
		return nil, ErrServerClosed
	}

	// 从 UDP 连接读取一个数据报
	data := make([]byte, PACKAGE_LIMIT)
	n, addr, err := conn.ReadFromUDP(data)
	if err != nil {
		return nil, err
	}

	lServer.Printf("%v bytes read from %v", n, addr)

	// 解码接收的数据报
	var header struct{ AskHeader }

	decoder := gob.NewDecoder(bytes.NewReader(data))

	// 解码数据报 header 部分
	if err := decoder.Decode(&header); err != nil {
		return nil, err
	}

	// 根据 action 生成数据报接收对象
	var pack Package
	switch header.Action {
	case ACTION_LOGIN:
		pack = &LoginAsk{}
	case ACTION_SHUTDOWN:
		pack = &ShutdownAsk{}
	}

	decoder = gob.NewDecoder(bytes.NewReader(data))

	// 解码完整数据报
	if err := decoder.Decode(pack); err != nil {
		return nil, err
	}

	// 生成请求对象
	return &Request{
		server: s,
		addr:   addr,
		action: header.Action,
		pack:   pack,
	}, nil
}

// 请求结构体
type Request struct {
	server    *Server      // UDP 服务器对象
	addr      *net.UDPAddr // 请求方地址
	sessionId SessionId    // 请求方 Session 编号
	action    ActionCode   // 本次请求 Action
	pack      Package      // 本次请求 Package
}

// 生成错误信息响应包
func (r *Request) makeErrorResponse(err error) (Package, error) {
	conn := r.server.conn
	if conn == nil {
		return nil, ErrServerClosed
	}

	pack := struct{ AckHeader }{
		AckHeader: AckHeader{
			Action:    r.action,
			SessionId: r.sessionId,
			IsOk:      false,
			Error:     err.Error(),
		},
	}

	return &pack, nil
}

// 处理登录请求
func (r *Request) handleActionLogin() (Package, error) {
	ask, ok := r.pack.(*LoginAsk)
	if !ok {
		return nil, ErrInvalidPackage
	}

	lServer.Printf("New ACTION_LOGIN body received, account=%v, password=%v", ask.Account, ask.Password)

	// 生成新的 SessionId
	if id, err := uuid.NewUUID(); err != nil {
		return nil, err
	} else {
		// 设置为当前 SessionId
		r.sessionId = SessionId(id.String())
	}

	// 生成登录响应对象
	return &LoginAck{
		AckHeader: AckHeader{
			Action:    ACTION_LOGIN,
			SessionId: r.sessionId,
			IsOk:      true,
		},
		Welcome: fmt.Sprintf("Welcome %v", ask.Account),
	}, nil
}

// 处理关闭服务请求
func (r *Request) handleActionShutdown() (Package, error) {
	ask, ok := r.pack.(*ShutdownAsk)
	if !ok {
		return nil, ErrInvalidPackage
	}

	// 判断请求中是否携带正确的 SessionId
	if !r.server.hasSession(ask.SessionId) {
		return r.makeErrorResponse(ErrInvalidSessionId)
	}

	lServer.Printf("New ACTION_SHUTDOWN body received")

	// 返回关闭服务响应对象, 附带错误信息要求关闭服务端
	return &ShutdownAck{
		AckHeader: AckHeader{
			Action:    ACTION_SHUTDOWN,
			SessionId: r.sessionId,
			IsOk:      true,
		},
	}, ErrCloseServer
}

// 响应结构体, 用在响应发送的 channel 上
type Response struct {
	addr  *net.UDPAddr // 接收方远程地址
	pack  Package      // 待发送的数据对象
	close bool         // 是否关闭服务器
}
