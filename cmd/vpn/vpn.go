package main

import (
	"go-vpn/netcheck"
)

type args struct {
	lPort int
	rAddr string
	rPort int
	dAddr string
}

func main() {
	//args := args{}
	//flag.IntVar(&args.lPort, "lport", 18080, "local udp port")
	//flag.StringVar(&args.rAddr, "raddr", "127.0.0.1", "remote udp listen addr")
	//flag.IntVar(&args.rPort, "rport", 18081, "remote udp listen port")
	//flag.StringVar(&args.dAddr, "daddr", "127.0.0.1", "remote udp listen port")
	//flag.Parse()

	//e := engine.Engine{
	//	LPort: args.lPort,
	//	RAddr: args.rAddr,
	//	RPort: args.rPort,
	//	DAddr: args.dAddr,
	//}
	//
	//err := e.Run()
	//if err != nil {
	//	log.Fatal(err)
	//}
	netcheck.NatTypeTest()
}
