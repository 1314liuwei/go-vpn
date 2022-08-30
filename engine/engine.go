package engine

import (
	"go-vpn/conn"
	"go-vpn/conn/udp"
	"log"
	"time"
)

type Engine struct {
	Addr string
	Port int
	Mode conn.Mode
	conn conn.Conn
}

func (e *Engine) Run() error {

	switch e.Mode {
	case conn.Server:
		e.Server()
	case conn.Client:
		e.Client()
	}

	return nil
}

func (e *Engine) Server() {
	server, err := udp.New(e.Addr, e.Port, conn.Server)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("server start work at %s:%d...\n", e.Addr, e.Port)
	for {
		buff := make([]byte, 1024)
		n, addr, err := server.ReadFrom(buff)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(string(buff[:n]))
		time.Sleep(time.Second)
		_, err = server.WriteTo([]byte("Hello, I'm server"), addr)
		if err != nil {
			panic(err)
		}
	}
}
func (e *Engine) Client() {
	client, err := udp.New(e.Addr, e.Port, conn.Client)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("client start work at %s:%d...\n", e.Addr, e.Port)
	for {
		_, err := client.Write([]byte("Hello, I'm client"))
		if err != nil {
			panic(err)
		}
		buff := make([]byte, 1024)
		n, err := client.Read(buff)
		if err != nil {
			log.Fatal(err)
		}
		time.Sleep(time.Second)
		log.Println(string(buff[:n]))
	}
}
