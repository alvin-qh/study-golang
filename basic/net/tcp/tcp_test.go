package tcp

import (
	"encoding/gob"
	"net"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func TestTcpConnect(t *testing.T) {
	server, err := StartTCPServer("0.0.0.0:8888")
	assert.NoError(t, err)
	defer server.StopTCPServer()

	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:8888")
	assert.NoError(t, err)

	conn, err := net.DialTCP("tcp", nil, addr)
	assert.NoError(t, err)
	defer conn.Close()

	conn.SetReadDeadline(time.Now().Add(time.Second * 10))
	conn.SetWriteDeadline(time.Now().Add(time.Second * 30))

	encoder := gob.NewEncoder(conn)
	decoder := gob.NewDecoder(conn)

	askHeader := TCPAskHeader{
		Action: ACTION_LOGIN,
	}
	err = encoder.Encode(&askHeader)
	assert.NoError(t, err)

	loginAsk := TCPLoginAsk{
		Account:  "Alvin",
		Password: "password",
	}
	err = encoder.Encode(&loginAsk)
	assert.NoError(t, err)

	ackHeader := TCPAckHeader{}
	err = decoder.Decode(&ackHeader)
	assert.NoError(t, err)

	assert.True(t, ackHeader.IsOk)

	loginAck := TCPLoginAck{}
	err = decoder.Decode(&loginAck)
	assert.NoError(t, err)

	encoder = gob.NewEncoder(conn)
	decoder = gob.NewDecoder(conn)

	askHeader = TCPAskHeader{
		Action: ACTION_SHUTDOWN,
	}
	err = encoder.Encode(&askHeader)
	assert.NoError(t, err)

	shutdownAsk := TCPShutdownAsk{}
	err = encoder.Encode(&shutdownAsk)
	assert.NoError(t, err)

	ackHeader = TCPAckHeader{}
	err = decoder.Decode(&ackHeader)
	assert.NoError(t, err)

	shutdownAck := TCPShutdownAck{}
	err = decoder.Decode(&shutdownAck)
	assert.NoError(t, err)

	server.Join()
}
