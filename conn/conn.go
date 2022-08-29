package conn

import (
	"io"
	"net"
)

type Mode int

func (m Mode) String() string {
	return []string{"Client", "Server"}[m]
}

const (
	Client Mode = iota
	Server
)

type Conn interface {
	io.ReadWriteCloser
	net.PacketConn
}
