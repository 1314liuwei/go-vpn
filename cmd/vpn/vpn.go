package main

import (
	"flag"
	"go-vpn/conn"
	"go-vpn/engine"
	"strings"
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

	var mode conn.Mode
	switch strings.ToLower(args.mode) {
	case "server":
		mode = conn.Server
	case "client":
		mode = conn.Client
	default:
		mode = conn.Client
	}

	e := engine.Engine{
		Mode: mode,
		Addr: args.addr,
		Port: args.port,
	}

	err := e.Run()
	if err != nil {
		return
	}
}
