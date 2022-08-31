package engine

import (
	"go-vpn/conn/udp"
	"log"
	"net/netip"
	"sync"
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
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		log.Println("start...")
		for {
			buff := []byte("hello world")
			addr, err := netip.ParseAddr(e.RAddr)
			if err != nil {
				panic(err)
			}
			_, err = e.conn.WriteToUDPAddrPort(buff, netip.AddrPortFrom(addr, uint16(e.RPort)))
			if err != nil {
				panic(err)
			}
			log.Println("send success!")
			time.Sleep(time.Second)
		}
		wg.Done()
	}()

	go func() {
		log.Println("start...")
		for {
			buff := make([]byte, 1024)
			n, endpoint, err := e.conn.ReadFromUDPAddrPort(buff)
			if err != nil {
				panic(err)
			}
			log.Printf("%s:%d >> %s", endpoint.Addr(), endpoint.Port(), buff[:n])
		}
		wg.Done()
	}()
	wg.Wait()
	return nil
}

func (e *Engine) ReceiveIPv4() {

}
