package main

import (
	"flag"
	"go-vpn/engine"
	"log"
)

type args struct {
	lPort int
	rAddr string
	rPort int
}

func main() {
	args := args{}
	flag.IntVar(&args.lPort, "lport", 18080, "local udp port")
	flag.StringVar(&args.rAddr, "raddr", "127.0.0.1", "remote udp listen addr")
	flag.IntVar(&args.rPort, "rport", 18081, "remote udp listen port")
	flag.Parse()

	e := engine.Engine{
		LPort: args.lPort,
		RAddr: args.rAddr,
		RPort: args.rPort,
	}

	err := e.Run()
	if err != nil {
		log.Fatal(err)
	}
}
