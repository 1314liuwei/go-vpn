package main

import (
	"go-vpn/config"
	"go-vpn/conn"
	"go-vpn/conn/udp"
	"log"

	"github.com/gogf/gf/v2/frame/g"

	"github.com/spf13/cobra"
)

var server = &cobra.Command{
	Use:   "server",
	Short: "",
	Run: func(cmd *cobra.Command, args []string) {
		conf, err := config.New("C:\\Users\\lw\\Desktop\\go-vpn\\config\\server.yaml")
		if err != nil {
			panic(err)
		}
		g.Dump(conf)

		server, err := udp.New(conf.Conn.Addr, conf.Conn.Port, conn.Server)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("server start work at %s:%d...\n", conf.Conn.Addr, conf.Conn.Port)
		for {
			buff := make([]byte, 1024)
			n, addr, err := server.ReadFrom(buff)
			if err != nil {
				log.Fatal(err)
			}
			log.Println(string(buff[:n]))
			_, err = server.WriteTo([]byte("Hello, I'm server"), addr)
			if err != nil {
				panic(err)
			}
		}
	},
}
