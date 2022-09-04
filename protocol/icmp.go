package protocol

import "go-vpn/util"

type ICMPRaw struct {
	Type     byte
	Code     byte
	CheckSum [2]byte
	Data     [1024]byte
}

var (
	ICMPControlMessages = map[int]map[int]string{
		0: {
			0: "Echo reply (used to ping)",
		},
	}
)

type ICMP struct {
	Type           int
	Code           int
	ControlMessage string // by Type and Code
	CheckSum       int
	Data           []byte
}

func (p Protocol) parseICMPRawPacket(packet ICMPRaw) *ICMP {
	icmp := &ICMP{}

	icmp.Type = util.Binary2Decimal(util.Bytes2Bits(packet.Type))
	icmp.Code = util.Binary2Decimal(util.Bytes2Bits(packet.Code))
	icmp.ControlMessage = parseICMPControlMessages(icmp.Type, icmp.Code)
	icmp.CheckSum = util.Binary2Decimal(util.Bytes2Bits(packet.CheckSum[0], packet.CheckSum[1]))
	icmp.Data = packet.Data[:]

	return icmp
}

func parseICMPControlMessages(t, c int) string {
	if codes, ok := ICMPControlMessages[t]; ok {
		if msg, ok := codes[c]; ok {
			return msg
		}
	}
	return "unknown message"
}
