package protocol

import (
	"bytes"
)

type MacRaw struct {
	Type [4]byte
}

type Type int

const (
	UnknownMacType Type = iota
	IP
)

func (t Type) String() string {
	return []string{"unKnownMacType", "IP"}[t]
}

type Mac struct {
	Type Type
}

func parseMacRawPacket(packet MacRaw) *Mac {
	mac := &Mac{}

	buff := make([]byte, len(packet.Type))
	copy(buff, packet.Type[:])
	if bytes.Equal(buff, []byte{0, 0, 8, 0}) {
		mac.Type = IP
	} else {
		mac.Type = UnknownMacType
	}

	return mac
}
