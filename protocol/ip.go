package protocol

import (
	"go-vpn/util"
	"net"
)

type IPProtocol int

const (
	ICMPProtocol IPProtocol = 1
	TCPProtocol  IPProtocol = 6
	UDPPProtocol IPProtocol = 17
)

type IPRaw struct {
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

type IP struct {
	Version        int // IP 协议版本号, IPv4 or IPv6
	HeaderLength   int // IP 报文首部长度
	DSCP           int
	TotalLength    int
	ID             int
	DF             bool
	MF             bool
	Offset         int
	TTL            int
	Protocol       IPProtocol
	HeaderChecksum int
	SourceAddr     net.IP
	DestAddr       net.IP
}

func (p Protocol) parseIPRawPacket(packet IPRaw) *IP {
	ip := &IP{}

	ip.Version, ip.HeaderLength = parseVersionAndHeaderLength(packet.VersionAndHeaderLength)
	ip.DSCP = util.Binary2Decimal(util.Bytes2Bits(packet.DifferentiatedServicesField)[:6])
	ip.TotalLength = parse2ByteToInt(packet.TotalLength)
	ip.ID = parse2ByteToInt(packet.Identification)
	ip.DF, ip.MF, ip.Offset = parseFlags(packet.Flags)
	ip.TTL = util.Binary2Decimal(util.Bytes2Bits(packet.TimeToLive))
	ip.Protocol = IPProtocol(util.Binary2Decimal(util.Bytes2Bits(packet.Protocol)))
	ip.HeaderChecksum = parse2ByteToInt(packet.HeaderChecksum)
	ip.SourceAddr = net.IPv4(packet.SourceAddr[0], packet.SourceAddr[1], packet.SourceAddr[2], packet.SourceAddr[3])
	ip.DestAddr = net.IPv4(packet.DestAddr[0], packet.DestAddr[1], packet.DestAddr[2], packet.DestAddr[3])

	return ip
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

func parse2ByteToInt(value [2]byte) int {
	buff := util.Bytes2Bits(value[0])
	buff = append(buff, util.Bytes2Bits(value[1])...)
	return util.Binary2Decimal(buff)
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
