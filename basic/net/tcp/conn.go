package tcp

import (
	"encoding/gob"
	"net"
	"sync/atomic"
)

// 包装 TCP 连接实例
type TCPConn struct {
	conn   *net.TCPConn // TCP 连接实例
	dec    *gob.Decoder // 接收数据的解码实例
	enc    *gob.Encoder // 发送数据的编码实例
	closed int32        // 连接是否已经关闭
}

// 创建实例
func NewTCPConn(conn *net.TCPConn) *TCPConn {
	return &TCPConn{
		conn:   conn,
		dec:    gob.NewDecoder(conn),
		enc:    gob.NewEncoder(conn),
		closed: 0,
	}
}

// 将数据编码后发送
func (c *TCPConn) Encode(input interface{}) error {
	return c.enc.Encode(input)
}

// 接收数据并解码
func (c *TCPConn) Decode(output interface{}) error {
	return c.dec.Decode(output)
}

// 获取远端连接地址
func (c *TCPConn) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

// 关闭当前连接
func (c *TCPConn) Close() error {
	if atomic.SwapInt32(&c.closed, 1) == 0 {
		return c.conn.Close()
	}
	return nil
}

// 返回当前连接是否已被关闭
func (c *TCPConn) IsClosed() bool {
	return atomic.LoadInt32(&c.closed) == 0
}
