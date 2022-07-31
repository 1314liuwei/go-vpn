package conn

import "io"

const (
	Server = iota
	Client
)

type Conn interface {
	io.ReadWriteCloser
}
