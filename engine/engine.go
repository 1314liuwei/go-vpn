package engine

import (
	"go-vpn/conn"
	"go-vpn/conn/udp"
)

type Engine struct {
	Addr string
	Port int
	conn conn.Conn
}

func (e *Engine) Run() error {
	var err error
	e.conn, err = udp.New(e.Port)
	if err != nil {
		return err
	}

	return nil
}

func (e *Engine) ReceiveIPv4() {

}
