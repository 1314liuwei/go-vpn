package protocol

import (
	"go-vpn/util"
	"net"
)

type IPv4Protocol int

//https://en.wikipedia.org/wiki/List_of_IP_protocol_numbers
const (
	ICMPProtocol IPv4Protocol = 1
	TCPProtocol  IPv4Protocol = 6
	UDPPProtocol IPv4Protocol = 17
)

type IPv4Raw struct {
	VersionAndHeaderLength      byte
	DifferentiatedServicesField byte
	TotalLength                 [2]byte
	Identification              [2]byte
	Flags                       [2]byte
	TimeToLive                  byte
	Protocol                    byte
	HeaderChecksum              [2]byte
	SourceAddr                  [4]byte
	DestAddr                    [4]byte
}

// IPv4 : https://en.wikipedia.org/wiki/IPv4
type IPv4 struct {
	Version        int // IPv4 协议版本号, IPv4 or IPv6
	HeaderLength   int // IPv4 报文首部长度
	DSCP           int
	TotalLength    int
	ID             int
	DF             bool
	MF             bool
	Offset         int
	TTL            int
	Protocol       IPv4Protocol
	HeaderChecksum int
	SourceAddr     net.IP
	DestAddr       net.IP
}

func (p Protocol) parseIPv4RawPacket(packet IPv4Raw) *IPv4 {
	ipv4 := &IPv4{}

	ipv4.Version, ipv4.HeaderLength = parseVersionAndHeaderLength(packet.VersionAndHeaderLength)
	ipv4.DSCP = util.Binary2Decimal(util.Bytes2Bits(packet.DifferentiatedServicesField)[:6])
	ipv4.TotalLength = util.Binary2Decimal(util.Bytes2Bits(packet.TotalLength[0], packet.TotalLength[1]))
	ipv4.ID = util.Binary2Decimal(util.Bytes2Bits(packet.Identification[0], packet.Identification[1]))
	ipv4.DF, ipv4.MF, ipv4.Offset = parseFlags(packet.Flags)
	ipv4.TTL = util.Binary2Decimal(util.Bytes2Bits(packet.TimeToLive))
	ipv4.Protocol = IPv4Protocol(util.Binary2Decimal(util.Bytes2Bits(packet.Protocol)))
	ipv4.HeaderChecksum = util.Binary2Decimal(util.Bytes2Bits(packet.HeaderChecksum[0], packet.HeaderChecksum[1]))
	ipv4.SourceAddr = net.IPv4(packet.SourceAddr[0], packet.SourceAddr[1], packet.SourceAddr[2], packet.SourceAddr[3])
	ipv4.DestAddr = net.IPv4(packet.DestAddr[0], packet.DestAddr[1], packet.DestAddr[2], packet.DestAddr[3])

	return ipv4
}

func parseVersionAndHeaderLength(value byte) (int, int) {
	var (
		version, headerLength int
	)
	buff := util.Bytes2Bits(value)
	version = util.Binary2Decimal(buff[:4])
	headerLength = util.Binary2Decimal(buff[4:])
	return version, headerLength
}

func parseFlags(value [2]byte) (bool, bool, int) {
	var (
		df, mf bool
		offset int
	)

	buff := util.Bytes2Bits(value[0])
	buff = append(buff, util.Bytes2Bits(value[1])...)
	df = buff[1] == 1
	mf = buff[2] == 1
	offset = util.Binary2Decimal(buff[3:])
	return df, mf, offset
}
