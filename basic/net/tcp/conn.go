package tcp

import (
	"encoding/gob"
	"net"
	"sync/atomic"
)

type Connection struct {
	conn   *net.TCPConn
	dec    *gob.Decoder
	enc    *gob.Encoder
	closed int32
}

func NewConnection(conn *net.TCPConn) *Connection {
	return &Connection{
		conn:   conn,
		dec:    gob.NewDecoder(conn),
		enc:    gob.NewEncoder(conn),
		closed: 0,
	}
}

func ConnectTo(addr *net.TCPAddr) (*Connection, error) {
	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		return nil, err
	}

	return &Connection{
		conn: conn,
		dec:  gob.NewDecoder(conn),
		enc:  gob.NewEncoder(conn),
	}, nil
}

func (c *Connection) Encode(input interface{}) error {
	return c.enc.Encode(input)
}

func (c *Connection) Decode(output interface{}) error {
	return c.dec.Decode(output)
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

func (c *Connection) Close() error {
	if atomic.SwapInt32(&c.closed, 1) == 0 {
		return c.conn.Close()
	}
	return nil
}

func (c *Connection) IsClosed() bool {
	return atomic.LoadInt32(&c.closed) == 0
}
