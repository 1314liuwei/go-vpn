package protocol

import (
	"bytes"
)

type MacRaw struct {
	Type [4]byte
}

// MacType : Mac 层承载的上层协议，如 IP 协议
type MacType int

const (
	UnknownMacType MacType = iota
	IPType
)

func (t MacType) String() string {
	return []string{"unKnownMacType", "IPType"}[t]
}

type Mac struct {
	Type MacType
}

func (p Protocol) parseMacRawPacket(packet MacRaw) *Mac {
	mac := &Mac{}

	buff := make([]byte, len(packet.Type))
	copy(buff, packet.Type[:])
	if bytes.Equal(buff, []byte{0, 0, 8, 0}) {
		mac.Type = IPType
	} else {
		mac.Type = UnknownMacType
	}

	return mac
}
