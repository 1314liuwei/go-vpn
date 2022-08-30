package udp

import (
	"fmt"
	"go-vpn/conn"
	"net"
	"time"
)

var _ conn.Conn = new(Connection)

type Connection struct {
	conn *net.UDPConn
}

func (c Connection) ReadFrom(p []byte) (n int, addr net.Addr, err error) {
	return c.conn.ReadFromUDP(p)
}

func (c Connection) WriteTo(p []byte, addr net.Addr) (n int, err error) {
	udpAddr, err := net.ResolveUDPAddr(addr.Network(), addr.String())
	if err != nil {
		return 0, err
	}
	return c.conn.WriteToUDP(p, udpAddr)
}

func (c Connection) LocalAddr() net.Addr {
	//TODO implement me
	panic("implement me")
}

func (c Connection) SetDeadline(t time.Time) error {
	//TODO implement me
	panic("implement me")
}

func (c Connection) SetReadDeadline(t time.Time) error {
	//TODO implement me
	panic("implement me")
}

func (c Connection) SetWriteDeadline(t time.Time) error {
	//TODO implement me
	panic("implement me")
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

func New(addr string, port int, mode conn.Mode) (conn.Conn, error) {
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
