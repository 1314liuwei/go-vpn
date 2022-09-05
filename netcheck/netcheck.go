package netcheck

import (
	"log"
	"net"
)

var (
	STUNServers = []string{
		"stun.qq.com:3478",
		"stun.bige0.com:3478",
	}
)

func NatTypeTest() {
	dial, err := net.Dial("udp4", STUNServers[0])
	if err != nil {
		log.Fatal(err)
		return
	}
	defer dial.Close()

	for {
		req := buildRequestPacket()
		_, err := dial.Write(req)
		if err != nil {
			log.Fatal(err)
			return
		}
		log.Println("send success!")

		buff := make([]byte, 1024)
		read, err := dial.Read(buff)
		if err != nil {
			log.Fatal(err)
			return
		}
		log.Printf("receive: %v", buff[:read])
		return
	}
}
