package protocol

import (
	"unsafe"

	"github.com/gogf/gf/v2/frame/g"
)

func Parse(buff []byte) {
	packet := make([]byte, len(buff))
	copy(packet, buff)

	macRaw := *(*MacRaw)(unsafe.Pointer(&packet[0]))
	mac := parseMacRawPacket(macRaw)

	g.Dump(mac)
}
