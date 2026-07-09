package wire

import (
	"net"
)

type RequestMessage struct {
	Request Request[[]byte]
	Conn    net.Conn
}

var RequestsQueue chan RequestMessage
