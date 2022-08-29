package main

import (
	"flag"

	"github.com/gogf/gf/v2/frame/g"
)

type args struct {
	addr string
	port int
	mode string
}

func main() {
	args := args{}
	flag.IntVar(&args.port, "port", 0, "udp port")
	flag.StringVar(&args.addr, "addr", "127.0.0.1", "udp listen addr")
	flag.StringVar(&args.mode, "mode", "client", "udp work mode")
	flag.Parse()

	g.Dump(args)
}
