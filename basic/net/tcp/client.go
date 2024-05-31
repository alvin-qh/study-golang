package tcp

import (
	"fmt"
	"net"
)

// 客户端结构体
type Client struct {
	conn *Connection
}

// 连接服务端
func Connect(address string) (*Client, error) {
	// 解析字符串地址
	addr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return nil, err
	}

	// 连接服务端
	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		return nil, err
	}
	lClient.Printf("Connect to server %v", addr)

	// 设置连接的发送和接收超时
	// conn.SetReadDeadline(time.Now().Add(time.Second * 10))
	// conn.SetWriteDeadline(time.Now().Add(time.Second * 30))

	// 返回 Client 结构体
	return &Client{
		conn: NewConnection(conn),
	}, nil
}

// 关闭连接
func (c *Client) Close() error {
	lClient.Printf("Connection %v closed", c.conn.RemoteAddr())
	return c.conn.Close()
}

// 发送请求数据, 返回
func (c *Client) Request(action ActionCode, body interface{}) (interface{}, error) {
	if err := c.sendRequest(action, body); err != nil {
		return nil, err
	}
	return c.receiveResponse(action)
}

// 发送请求
func (c *Client) sendRequest(action ActionCode, body interface{}) error {
	// 发送请求头
	if err := c.conn.Encode(&AskHeader{Action: action}); err != nil {
		return err
	}
	lClient.Printf("Send ask header to %v, action=%v", c.conn.RemoteAddr(), action)

	// 发送请求内容
	if err := c.conn.Encode(body); err != nil {
		return err
	}
	lClient.Printf("Send ask body to %v, action=%v", c.conn.RemoteAddr(), action)

	return nil
}

// 接收响应
func (c *Client) receiveResponse(action ActionCode) (interface{}, error) {
	// 接收响应头
	header := AckHeader{}
	if err := c.conn.Decode(&header); err != nil {
		return nil, err
	}

	// 确认响应正确
	if header.Action != action {
		return nil, fmt.Errorf("invalid response action %v", header.Action)
	}
	lClient.Printf("Receive ack header from %v, action=%v", c.conn.RemoteAddr(), action)

	// 根据响应头接收响应内容
	var resp interface{} = nil
	switch action {
	case ACTION_LOGIN:
		resp = &LoginAck{}
	case ACTION_SHUTDOWN:
		resp = &ShutdownAck{}
	default:
		return nil, fmt.Errorf("invalid action code %v", action)
	}

	if err := c.conn.Decode(resp); err != nil {
		return nil, err
	}
	lClient.Printf("Receive ack body from %v, action=%v", c.conn.RemoteAddr(), action)

	return resp, nil
}
