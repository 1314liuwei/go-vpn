package engine

import (
	"go-vpn/conn"
	"go-vpn/conn/udp"
	"log"
	"net/netip"
	"time"
)

type Engine struct {
	LPort int
	RAddr string
	RPort int
	conn  *udp.Connection
}

func (e *Engine) Run() error {
	var err error
	e.conn, err = udp.New(e.LPort)
	if err != nil {
		return err
	}

	log.Println("start run...")

	addr, err := netip.ParseAddr(e.RAddr)
	if err != nil {
		panic(err)
	}
	e.conn.Notify(conn.Op{Action: "add", Value: netip.AddrPortFrom(addr, uint16(e.RPort))})

	go SendIPv4(e.conn)
	go ReceiveIPv4(e.conn)
	select {}
	return nil
}

func ReceiveIPv4(conn *udp.Connection) {
	for {
		buff := make([]byte, 1024)
		n, endpoint, err := conn.ReadFromUDPAddrPort(buff)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("%s:%d >> %s", endpoint.Addr(), endpoint.Port(), buff[:n])
	}
}

func SendIPv4(conn *udp.Connection) {
	log.Println("start send...")
	for {
		buff := []byte("hello world")
		_, err := conn.Write(buff)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("send success!")
		time.Sleep(time.Second)
	}
}
