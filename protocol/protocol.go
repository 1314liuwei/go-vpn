package protocol

import (
	"unsafe"

	"github.com/gogf/gf/v2/frame/g"
)

type Protocol struct{}

func (p Protocol) Parse(buff []byte) {
	packet := make([]byte, len(buff))
	copy(packet, buff)

	macRaw := *(*MacRaw)(unsafe.Pointer(&packet[0]))
	mac := p.parseMacRawPacket(macRaw)
	macSize := unsafe.Sizeof(MacRaw{})
	g.Dump(mac)
	switch mac.Type {
	case IPType:
		//ipSize := unsafe.Sizeof(IP{})
		ipRaw := *(*IPRaw)(unsafe.Pointer(&packet[macSize]))
		ip := p.parseIPRawPacket(ipRaw)
		g.Dump(ip)
	default:
		return
	}
}

func New() *Protocol {
	return &Protocol{}
}
