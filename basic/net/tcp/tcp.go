package tcp

import (
	"encoding/gob"
	"log"
	"net"
	"os"
	"time"
)

var (
	logger = log.New(os.Stdout, "DEBUG ", log.LstdFlags|log.Lshortfile|log.Ltime)
)

type Server struct {
	listener *net.TCPListener
	closeCh  chan struct{}
}

func StartTCPServer(address string) (*Server, error) {
	addr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return nil, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		logger.Fatalf("network error: %v", err)
		return nil, err
	}
	server := &Server{
		listener: l,
		closeCh:  make(chan struct{}),
	}

	go server.handleAcception(l)

	return server, nil
}

func (s *Server) Join() {
	<-s.closeCh
}

func (s *Server) StopTCPServer() error {
	var err error = nil
	if s.listener != nil {
		err = s.listener.Close()
		s.listener = nil
	}
	return err
}

func (s *Server) handleAcception(l *net.TCPListener) {
	defer close(s.closeCh)
    defer s.listener.Close()

	logger.Printf("Ready to accepting, %v", l.Addr())

	for {
		conn, err := l.AcceptTCP()
		if err != nil {
			logger.Fatalf("Network error: %v", err)
			break
		}

		go s.handleClientSession(conn)
	}
}

type Request struct {
	server  *Server
	conn    *net.TCPConn
	decoder *gob.Decoder
	encoder *gob.Encoder
	header  TCPAskHeader
	context map[string]interface{}
}

func (r *Request) decodeAskHeader() (TCPAskHeader, error) {
	if err := r.decoder.Decode(&r.header); err != nil {
		logger.Printf("Decode ask header failed: %v", err)
		return r.header, err
	}
	return r.header, nil
}

func (r *Request) decodeAskBody(body interface{}) error {
	if err := r.decoder.Decode(body); err != nil {
		logger.Printf("Decode ask body failed: %v", err)
		return err
	}
	return nil
}

func (r *Request) encodeAckHeader(header *TCPAckHeader) error {
	if err := r.encoder.Encode(header); err != nil {
		logger.Printf("Encode ack header failed: %v", err)
		return err
	}
	return nil
}

func (r *Request) encodeAckBody(body interface{}) error {
	if err := r.encoder.Encode(body); err != nil {
		logger.Printf("Encode ack body failed: %v", err)
	}
	return nil
}

// cspell: ignore sess
func (s *Server) handleClientSession(conn *net.TCPConn) {
	logger.Printf("New connection comming, %v", conn.RemoteAddr())
	defer conn.Close()

	conn.SetKeepAlive(true)
	for {
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

func (r *Request) handleLoginAction() error {
	ask := TCPLoginAsk{}
	if err := r.decodeAskBody(&ask); err != nil {
		return err
	}
	logger.Printf("New ACTION_LOGIN body received, account=%v, password=%v", ask.Account, ask.Password)

	header := TCPAckHeader{
		Action: ACTION_LOGIN,
		IsOk:   true,
	}
	if err := r.encodeAckHeader(&header); err != nil {
		return err
	}

	body := TCPLoginAck{}
	if err := r.encodeAckBody(&body); err != nil {
		return err
	}

	return nil
}

func (r *Request) handleShutdownAction() error {
	ask := TCPShutdownAsk{}
	if err := r.decodeAskBody(&ask); err != nil {
		return err
	}
	logger.Printf("New ACTION_SHUTDOWN body received")

	header := TCPAckHeader{
		Action: ACTION_SHUTDOWN,
		IsOk:   true,
	}
	if err := r.encodeAckHeader(&header); err != nil {
		return err
	}

	body := TCPShutdownAsk{}
	if err := r.encodeAckBody(&body); err != nil {
		return err
	}

	r.server.StopTCPServer()
	return nil
}
