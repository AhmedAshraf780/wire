package wire

import "net"

type ResponseMessage struct {
	Conn     net.Conn
	Response Response
}

var ResponseQueue chan ResponseMessage
