package udp

import (
	"fmt"
	"go-vpn/conn"
	"net"
)

var _ conn.Conn = new(Connection)

type Connection struct {
	conn *net.UDPConn
}

func (c Connection) Read(p []byte) (n int, err error) {
	return c.conn.Read(p)
}

func (c Connection) Write(p []byte) (n int, err error) {
	return c.conn.Write(p)
}

func (c Connection) Close() error {
	return c.conn.Close()
}

func New(addr string, port int, mode int) (conn.Conn, error) {
	switch mode {
	case conn.Server:
		udp, err := net.ListenUDP("udp", &net.UDPAddr{
			IP:   net.ParseIP(addr).To4(),
			Port: port,
		})
		if err != nil {
			return nil, err
		}

		return &Connection{
			conn: udp,
		}, nil
	case conn.Client:
		udp, err := net.DialUDP("udp", nil, &net.UDPAddr{
			IP:   net.ParseIP(addr),
			Port: port,
		})
		if err != nil {
			return nil, err
		}

		return &Connection{
			conn: udp,
		}, nil
	default:
		return nil, fmt.Errorf("unknow running mode")
	}

}
