package main

import (
	"go-vpn/conn"
	"go-vpn/conn/udp"
	"log"

	"github.com/spf13/cobra"
)

var server = &cobra.Command{
	Use:   "server",
	Short: "",
	Run: func(cmd *cobra.Command, args []string) {
		server, err := udp.New("127.0.0.1", 8088, conn.Server)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("server start work at %s:%d...\n", "127.0.0.1", 8088)
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
