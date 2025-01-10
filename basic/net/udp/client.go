package udp

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net"
)

// 客户端结构体
type Client struct {
	conn      *net.UDPConn
	addr      *net.UDPAddr
	sessionId SessionId
}

// 连接服务端
func Connect(address string) (*Client, error) {
	// 解析字符串地址
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return nil, err
	}

	// 连接服务端
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return nil, err
	}

	// 设置连接的发送和接收超时
	// conn.SetReadDeadline(time.Now().Add(time.Second * 10))
	// conn.SetWriteDeadline(time.Now().Add(time.Second * 30))

	// 返回 Client 结构体
	return &Client{
		conn: conn,
		addr: addr,
	}, nil
}

// 关闭连接
func (c *Client) Close() error {
	var err error
	if c.conn != nil {
		// 关闭与服务端的连接
		err = c.conn.Close()
		c.conn = nil
	}
	return err
}

// 发送请求数据, 返回
func (c *Client) Request(pack Package) (Package, error) {
	if err := c.sendRequest(pack); err != nil {
		return nil, err
	}
	return c.receiveResponse(pack.GetAction())
}

// 发送请求
func (c *Client) sendRequest(pack Package) error {
	buf := bytes.NewBuffer(make([]byte, 0, 1024))
	encoder := gob.NewEncoder(buf)

	pack.SetSessionId(c.sessionId)

	// 发送请求头
	if err := encoder.Encode(pack); err != nil {
		return err
	}

	// 发送数据
	n, err := c.conn.Write(buf.Bytes())
	if err != nil {
		return err
	}

	cLog.Printf("%v bytes was sent to %v", n, c.addr)
	return nil
}

// 接收响应
func (c *Client) receiveResponse(action ActionCode) (Package, error) {
	// 接收响应数据
	data := make([]byte, 1024)
	n, addr, err := c.conn.ReadFromUDP(data)
	if err != nil {
		return nil, err
	}
	cLog.Printf("%v bytes was received from %v", n, addr)

	decoder := gob.NewDecoder(bytes.NewReader(data))

	// 解码响应数据头
	var header struct{ AckHeader }
	if err := decoder.Decode(&header); err != nil {
		return nil, err
	}

	// 确认响应正确
	if header.Action != action {
		return nil, fmt.Errorf("invalid response action %v", header.Action)
	}

	cLog.Printf("received session id: %v, action=%v", header.SessionId, action)

	// 确认是否携带错误信息
	if !header.IsOk {
		return nil, fmt.Errorf("%v", header.Error)
	}

	// 保存 session id
	c.sessionId = header.SessionId

	// 根据响应类型创建响应数据类型
	var resp Package = nil
	switch action {
	case ACTION_LOGIN:
		resp = &LoginAck{}
	case ACTION_SHUTDOWN:
		resp = &ShutdownAck{}
	default:
		return nil, fmt.Errorf("invalid action code %v", action)
	}

	// 解码完整的数据包
	decoder = gob.NewDecoder(bytes.NewReader(data))
	if err := decoder.Decode(resp); err != nil {
		return nil, err
	}

	return resp, nil
}
