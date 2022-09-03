package protocol

import (
	"unsafe"

	"github.com/gogf/gf/v2/frame/g"
)

type Protocol struct{}

/*
Parse 报文解析器
根据协议解析报文
*/
func (p Protocol) Parse(buff []byte) {
	packet := make([]byte, len(buff))
	copy(packet, buff)

	macRaw := *(*MacRaw)(unsafe.Pointer(&packet[0]))
	mac := p.parseMacRawPacket(macRaw)
	macSize := unsafe.Sizeof(MacRaw{})
	g.Dump(mac)
	switch mac.Type {
	case IPv4Type:
		//ipSize := unsafe.Sizeof(IPv4{})
		ipRaw := *(*IPv4Raw)(unsafe.Pointer(&packet[macSize]))
		ip := p.parseIPv4RawPacket(ipRaw)
		g.Dump(ip)
	default:
		return
	}
}

func Parse(packet []byte) {
	p := Protocol{}
	p.Parse(packet)
	return
}
