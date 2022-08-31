package udp

import (
	"go-vpn/conn"
	"net"
	"net/netip"
	"time"
)

var _ conn.Conn = new(Connection)

type Connection struct {
	ch        chan conn.Op
	endpoints []netip.AddrPort
	conn      *net.UDPConn
}

func (c *Connection) ReadFrom(p []byte) (n int, addr net.Addr, err error) {
	return c.conn.ReadFromUDP(p)
}

func (c *Connection) WriteTo(p []byte, addr net.Addr) (n int, err error) {
	udpAddr, err := net.ResolveUDPAddr(addr.Network(), addr.String())
	if err != nil {
		return 0, err
	}
	return c.conn.WriteToUDP(p, udpAddr)
}

func (c *Connection) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}

func (c *Connection) SetDeadline(t time.Time) error {
	//TODO implement me
	panic("implement me")
}

func (c *Connection) SetReadDeadline(t time.Time) error {
	//TODO implement me
	panic("implement me")
}

func (c *Connection) SetWriteDeadline(t time.Time) error {
	//TODO implement me
	panic("implement me")
}

func (c *Connection) Read(p []byte) (n int, err error) {
	return c.conn.Read(p)
}

func (c *Connection) Write(p []byte) (n int, err error) {
	for _, endpoint := range c.endpoints {
		_, err := c.conn.WriteToUDPAddrPort(p, endpoint)
		if err != nil {
			return 0, err
		}
	}
	return
}

func (c *Connection) WriteToUDPAddrPort(b []byte, addr netip.AddrPort) (int, error) {
	return c.conn.WriteToUDPAddrPort(b, addr)
}

func (c *Connection) ReadFromUDPAddrPort(b []byte) (n int, addr netip.AddrPort, err error) {
	return c.conn.ReadFromUDPAddrPort(b)
}

func (c *Connection) Close() error {
	return c.conn.Close()
}

func (c *Connection) Notify(op conn.Op) {
	switch op.Action {
	case "add":
		endpoint := op.Value.(netip.AddrPort)
		c.endpoints = append(c.endpoints, endpoint)
	}
}

func New(port int) (*Connection, error) {
	udp, err := net.ListenUDP("udp4", &net.UDPAddr{
		Port: port,
	})
	if err != nil {
		return nil, err
	}

	return &Connection{
		conn: udp,
	}, nil
}
