package engine

import (
	"go-vpn/conn/udp"
	"log"
	"net/netip"
	"sync"
)

type Engine struct {
	Addr string
	Port int
	conn *udp.Connection
}

func (e *Engine) Run() error {
	var err error
	e.conn, err = udp.New(e.Port)
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
			addr, err := netip.ParseAddr("127.0.0.1")
			if err != nil {
				log.Fatal("parse: ", err)
			}
			_, err = e.conn.WriteToUDPAddrPort(buff, netip.AddrPortFrom(addr, 8088))
			if err != nil {
				log.Fatal("w:", err)
			}
		}
		wg.Done()
	}()

	go func() {
		log.Println("start...")
		for {
			buff := make([]byte, 1024)
			_, endpoint, err := e.conn.ReadFromUDPAddrPort(buff)
			if err != nil {
				panic(err)
			}
			log.Println(endpoint)
		}
		wg.Done()
	}()
	wg.Wait()
	return nil
}

func (e *Engine) ReceiveIPv4() {

}
