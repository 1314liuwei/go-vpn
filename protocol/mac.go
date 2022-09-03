package protocol

import (
	"bytes"
)

type MacRaw struct {
	Type [4]byte
}

// EtherType : Mac 层承载的上层协议，如 IPv4 协议
type EtherType int

// https://en.wikipedia.org/wiki/EtherType
const (
	UnknownEtherType EtherType = iota
	IPv4Type
	ARP
)

func (t EtherType) String() string {
	return []string{"unKnownMacType", "IPv4Type"}[t]
}

type Mac struct {
	Type EtherType
}

func (p Protocol) parseMacRawPacket(packet MacRaw) *Mac {
	mac := &Mac{}

	buff := make([]byte, len(packet.Type))
	copy(buff, packet.Type[:])

	// 大端小端
	if bytes.Equal(buff, []byte{0, 0, 8, 0}) {
		mac.Type = IPv4Type
	} else if bytes.Equal(buff, []byte{6, 0, 8, 0}) {
		mac.Type = ARP
	} else {
		mac.Type = UnknownEtherType
	}

	return mac
}
