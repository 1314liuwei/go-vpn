package main

import (
	"go-vpn/conn"
	"go-vpn/conn/udp"
	"log"

	"github.com/spf13/cobra"
)

var client = &cobra.Command{
	Use:   "client",
	Short: "",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := udp.New("127.0.0.1", 8088, conn.Client)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("client start work at %s:%d...\n", "127.0.0.1", 8088)
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
			log.Println(string(buff[:n]))
		}
	},
}
