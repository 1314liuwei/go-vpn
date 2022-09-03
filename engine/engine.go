package engine

import (
	"go-vpn/conn"
	"go-vpn/conn/udp"
	"go-vpn/device"
	"go-vpn/protocol"
	"log"
	"net/netip"
)

type Engine struct {
	LPort  int
	RAddr  string
	RPort  int
	DAddr  string
	conn   *udp.Connection
	device device.Manager
}

func (e *Engine) Run() error {
	var err error
	e.conn, err = udp.New(e.LPort)
	if err != nil {
		return err
	}

	e.AddEndpoint(e.RAddr)
	e.CreateDevice(e.DAddr)
	go SendIPv4(e.conn, e.device)
	go ReceiveIPv4(e.conn, e.device)
	select {}
}

func (e *Engine) CreateDevice(addr string) {
	name := "tun0"
	tun, err := device.CreateTUN(name)
	if err != nil {
		log.Fatal(err)
	}
	//defer tun.Close()

	wrap := device.NewManage(tun)
	wrap.SetMTU(1280)
	err = wrap.ChangeState(device.UP)
	if err != nil {
		log.Fatal(err)
	}

	err = wrap.SetAddrIPv4(addr)
	if err != nil {
		log.Fatal(err)
	}
	e.device = wrap
}

func (e *Engine) AddEndpoint(addr string) {
	add, err := netip.ParseAddr(addr)
	if err != nil {
		panic(err)
	}
	e.conn.Notify(conn.Op{Action: "add", Value: netip.AddrPortFrom(add, uint16(e.RPort))})
}

func ReceiveIPv4(conn *udp.Connection, dev device.Manager) {
	log.Println("start receive...")
	for {
		buff := make([]byte, 1024)
		n, endpoint, err := conn.ReadFromUDPAddrPort(buff)
		if err != nil {
			log.Fatal(err)
		}

		_, err = dev.Write(buff[:n])
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("%s:%d >> %v", endpoint.Addr(), endpoint.Port(), buff[:n])
	}
}

func SendIPv4(conn *udp.Connection, dev device.Manager) {
	log.Println("start send...")
	for {
		buff := make([]byte, 1024)
		n, err := dev.Read(buff)
		if err != nil {
			log.Fatal(err)
		}

		p := protocol.New()
		p.Parse(buff[:n])
		_, err = conn.Write(buff[:n])
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("send >> %v", buff[:n])
	}
}
