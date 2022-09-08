package netcheck

import (
	"encoding/binary"
	"go-vpn/util"
	"log"
	"math/rand"
	"net"
	"time"
	"unsafe"
)

const (
	STUNBindingRequest       = 0x0001
	STUNBindingResponse      = 0x0101
	STUNBindingErrorResponse = 0x0111

	MappedAddress     = 0x0001
	ResponseAddress   = 0x0002
	ChangeRequest     = 0x0003
	SourceAddress     = 0x0004
	ChangedAddress    = 0x0005
	USERNAME          = 0x0006
	PASSWORD          = 0x0007
	MessageIntegrity  = 0x0008
	ErrorCode         = 0x0009
	UnknownAttributes = 0x000a
	ReflectedFrom     = 0x000b
)

var (
	AttributesName = map[int]string{
		MappedAddress:     "MAPPED_ADDRESS",
		ResponseAddress:   "RESPONSE_ADDRESS",
		ChangeRequest:     "CHANGE_REQUEST",
		SourceAddress:     "SOURCE_ADDRESS",
		ChangedAddress:    "CHANGED_ADDRESS",
		USERNAME:          "USERNAME",
		PASSWORD:          "PASSWORD",
		MessageIntegrity:  "MESSAGE_INTEGRITY",
		ErrorCode:         "ERROR_CODE",
		UnknownAttributes: "UNKNOWN_ATTRIBUTES",
		ReflectedFrom:     "REFLECTED_FROM",
	}
)

type STUNHeaderPacket struct {
	MessageType   [2]byte
	MessageLength [2]byte
	MagicCookie   [4]byte
	TransactionID [12]byte
}

type AddressAttribute struct {
	Reserved       int
	ProtocolFamily int
	Port           int
	IP             net.IP
}

func buildRequestHeader(mLen int) []byte {
	req := &STUNHeaderPacket{}

	buff := make([]byte, 1024)
	binary.BigEndian.PutUint16(buff, STUNBindingRequest)
	copy(req.MessageType[:], buff)

	binary.BigEndian.PutUint16(buff, uint16(mLen))
	copy(req.MessageLength[:], buff)

	binary.BigEndian.PutUint32(buff, 0x2112A442)
	copy(req.MagicCookie[:], buff)

	rand.Seed(time.Now().Unix())
	tid := [12]byte{}
	for i := 0; i < 12; i++ {
		tid[i] = byte(rand.Int())
	}
	req.TransactionID = tid

	*(*STUNHeaderPacket)(unsafe.Pointer(&buff[0])) = *req
	size := unsafe.Sizeof(STUNHeaderPacket{})
	return buff[:size]
}

func buildChangePortAndIPRequest(ip, port bool) []byte {
	var packet []byte
	header := buildRequestHeader(8)
	packet = append(packet, header...)

	// Attribute Type
	buff := make([]byte, 1024)
	binary.BigEndian.PutUint16(buff, ChangeRequest)
	packet = append(packet, buff[:2]...)

	// Attribute Length
	binary.BigEndian.PutUint16(buff, 4)
	packet = append(packet, buff[:2]...)

	// Set IP and Port: 00000000 00000000 00000000 00000(IP)(Port)0
	set := 0
	if ip {
		set += 4
	}
	if port {
		set += 2
	}
	binary.BigEndian.PutUint32(buff, uint32(set))
	packet = append(packet, buff[:4]...)
	return packet
}

func parseResponseAttributes(buff []byte) (map[string]AddressAttribute, error) {
	var (
		i      int
		out    = map[string]AddressAttribute{}
		header STUNHeaderPacket
	)

	header = *(*STUNHeaderPacket)(unsafe.Pointer(&buff[0]))
	size := int(unsafe.Sizeof(STUNHeaderPacket{}))
	i = size

	log.Printf("header: %v\n", header)
	log.Printf("buff: %v\n", buff)

	for i-size < util.Binary2Decimal(util.Bytes2Bits(header.MessageLength[0], header.MessageLength[1])) {
		AttributeTypeValue := util.Binary2Decimal(util.Bytes2Bits(buff[i], buff[i+1]))
		i += 2
		AttributeLength := util.Binary2Decimal(util.Bytes2Bits(buff[i], buff[i+1]))
		i += 2

		switch AttributeTypeValue {
		case MappedAddress, SourceAddress, ChangedAddress:
			if name, ok := AttributesName[AttributeTypeValue]; ok {
				out[name] = parseAddressAttribute(buff[i : i+AttributeLength])
			}
		}
		i += AttributeLength
	}

	return out, nil
}

func parseAddressAttribute(buff []byte) AddressAttribute {
	Reserved := util.Binary2Decimal(util.Bytes2Bits(buff[0]))
	ProtocolFamily := util.Binary2Decimal(util.Bytes2Bits(buff[1]))
	Port := util.Binary2Decimal(util.Bytes2Bits(buff[2], buff[3]))
	IP := net.IPv4(buff[4], buff[5], buff[6], buff[7]).To4()

	out := AddressAttribute{
		Reserved:       Reserved,
		ProtocolFamily: ProtocolFamily,
		Port:           Port,
		IP:             IP,
	}
	return out
}
