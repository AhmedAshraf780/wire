package wire

import "net"

type ResponseMessage struct {
	Conn     net.Conn
	Response Response[[]byte]
}

var ResponseQueue chan ResponseMessage
