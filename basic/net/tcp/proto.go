package tcp

import "errors"

type ActionCode int

func (a ActionCode) String() string {
	switch a {
	case ACTION_LOGIN:
		return "ACTION_LOGIN"
	case ACTION_SHUTDOWN:
		return "ACTION_SHUTDOWN"
	default:
		panic(errors.New("invalid action code"))
	}
}

const (
	ACTION_LOGIN ActionCode = iota
	ACTION_SHUTDOWN
)

type TCPAskHeader struct {
	Action ActionCode
}

type TCPAckHeader struct {
	Action ActionCode
	IsOk   bool
	Error  string
}

type TCPLoginAsk struct {
	Account  string
	Password string
}

type TCPLoginAck struct{}

type TCPShutdownAsk struct{}

type TCPShutdownAck struct{}
