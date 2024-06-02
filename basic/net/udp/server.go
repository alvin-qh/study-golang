package udp

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"log"
	"net"
	"sync"
	"sync/atomic"
	"unsafe"

	"github.com/google/uuid"
)

// 定义错误值
var (
	ErrInvalidPackage    = errors.New("invalid package")
	ErrInvalidSessionId  = errors.New("invalid session id")
	ErrServerShouldClose = fmt.Errorf("server should close")
)

// 定义 Server 结构体
type Server struct {
	conn    *net.UDPConn   // UDP 连接对象
	session sync.Map       // 保存 session 信息
	wg      sync.WaitGroup // 发送接收 goroutine 等待组
}

// 响应结构体, 用在响应发送的 channel 上
type Response struct {
	addr *net.UDPAddr // 接收方远程地址
	pack Package      // 待发送的数据对象
}

// 启动服务器
func ServerStart(address string) (*Server, error) {
	// 解析监听地址
	addr, err := net.ResolveUDPAddr("udp4", address)
	if err != nil {
		sLog.Fatalf("Cannot resolve address %v", address)
		return nil, err
	}

	// 监听指定端口和地址
	conn, err := net.ListenUDP("udp4", addr)
	if err != nil {
		sLog.Fatalf("Cannot Listen UDP at %v", address)
		return nil, err
	}
	sLog.Printf("Server was listened at %v", addr)

	// 创建服务端对象
	srv := &Server{
		conn:    conn,
		session: sync.Map{},
	}

	// 设置两个等待任务
	srv.wg.Add(2)

	// 创建响应发送信道
	ch := make(chan Response, 100)

	// 接收信息
	go srv.handleReceiveMessage(ch)

	// 发送信息
	go srv.handleSendMessage(ch)

	// 处理未捕获异常
	defer func() {
		r := recover()
		if err, ok := r.(error); ok {
			sLog.Panicf("Panic %v", err)
		}
	}()

	return srv, nil
}

// 处理接收数据
func (s *Server) handleReceiveMessage(ch chan<- Response) {
	defer func() {
		close(ch)

		s.wg.Done()
		s.Close()
	}()

	// 获取 UDP 连接对象
	conn := s.conn
	if conn == nil {
		return
	}

	closed := false

	for !closed {
		// 接收数据
		req, err := receivePackage(conn)
		if err != nil {
			sLog.Printf("Cannot receive package, caused: %v", err)
			break
		}
		sLog.Printf("Received package from %v, action=%v, session-id=%v", req.addr, req.action, req.sessionId)

		// 处理已接收的数据
		go func() {
			var pack Package

			switch req.action {
			case ACTION_LOGIN:
				// 处理登录消息
				pack, err = req.handleActionLogin()

				// 保存 session 对象
				s.session.Store(req.sessionId, struct{}{})
			case ACTION_SHUTDOWN:
				// 处理关闭服务消息
				pack, err = req.handleActionShutdown(&s.session)
			}

			if err != nil {
				if errors.Is(err, ErrServerShouldClose) {
					closed = true
					return
				} else {
					log.Fatal(err)
				}
			}

			// 发送数据到 channel
			ch <- Response{
				addr: req.addr,
				pack: pack,
			}
		}()
	}
}

// 处理发送信息
func (s *Server) handleSendMessage(ch <-chan Response) {
	defer s.wg.Done()

	// 获取 UDP 连接对象
	conn := s.conn
	if conn == nil {
		return
	}

	// 从 channel 中获取待发送的数据
	for resp := range ch {
		// 将数据写入 UDP 连接
		buf := bytes.NewBuffer(make([]byte, 0, PACKAGE_LIMIT))
		enc := gob.NewEncoder(buf)

		// 将发送数据编码后写入缓冲
		if err := enc.Encode(resp.pack); err != nil {
			sLog.Printf("Cannot send response, caused %v", err)
		}

		// 缓冲数据写入 UDP
		if n, err := conn.WriteToUDP(buf.Bytes(), resp.addr); err != nil {
			sLog.Printf("Cannot send response, caused %v", err)
			break
		} else {
			sLog.Printf("%v bytes write to %v", n, resp.addr)
		}
	}
}

// 等待服务器关闭
func (s *Server) Close() {
	conn := (*net.UDPConn)(atomic.SwapPointer((*unsafe.Pointer)(unsafe.Pointer(&s.conn)), nil))
	if conn != nil {
		conn.Close()
		s.wg.Wait()
	}
}

// 接收数据包
func receivePackage(conn *net.UDPConn) (*Request, error) {
	// 从 UDP 连接读取一个数据报
	data := make([]byte, PACKAGE_LIMIT)
	n, addr, err := conn.ReadFromUDP(data)
	if err != nil {
		return nil, err
	}

	sLog.Printf("%v bytes read from %v", n, addr)

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
		addr:      addr,
		action:    header.Action,
		pack:      pack,
		sessionId: header.SessionId,
	}, nil
}

// 请求结构体
type Request struct {
	addr      *net.UDPAddr // 请求方地址
	sessionId SessionId    // 请求方 Session 编号
	action    ActionCode   // 本次请求 Action
	pack      Package      // 本次请求 Package
}

// 处理登录请求
func (r *Request) handleActionLogin() (Package, error) {
	ask, ok := r.pack.(*LoginAsk)
	if !ok {
		return nil, ErrInvalidPackage
	}

	sLog.Printf("New ACTION_LOGIN body received, account=%v, password=%v", ask.Account, ask.Password)

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
func (r *Request) handleActionShutdown(session *sync.Map) (Package, error) {
	ask, ok := r.pack.(*ShutdownAsk)
	if !ok {
		return nil, ErrInvalidPackage
	}

	// 判断请求中是否携带正确的 SessionId
	if _, ok := session.Load(ask.SessionId); !ok {
		return r.makeErrorResponse(ErrInvalidSessionId)
	}

	sLog.Printf("New ACTION_SHUTDOWN body received")

	// 返回关闭服务响应对象, 附带错误信息要求关闭服务端
	return &ShutdownAck{
		AckHeader: AckHeader{
			Action:    ACTION_SHUTDOWN,
			SessionId: r.sessionId,
			IsOk:      true,
		},
	}, nil
}

// 生成错误信息响应包
func (r *Request) makeErrorResponse(err error) (Package, error) {
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
