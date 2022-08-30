package main

import (
	"flag"
	"go-vpn/engine"
	"log"
)

type args struct {
	addr string
	port int
}

func main() {
	args := args{}
	flag.IntVar(&args.port, "port", 0, "udp port")
	flag.StringVar(&args.addr, "addr", "127.0.0.1", "udp listen addr")
	flag.Parse()

	e := engine.Engine{
		Addr: args.addr,
		Port: args.port,
	}

	err := e.Run()
	if err != nil {
		log.Fatal(err)
	}
}
