package wire

import (
	"net"
)

type RequestMessage struct {
	Request Request
	Conn    net.Conn
}

var RequestsQueue chan RequestMessage
