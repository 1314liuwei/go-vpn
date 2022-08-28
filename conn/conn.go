package conn

import (
	"io"
	"net"
)

const (
	Server = iota
	Client
)

type Conn interface {
	io.ReadWriteCloser
	net.PacketConn
}
