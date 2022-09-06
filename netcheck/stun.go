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

	MAPPED_ADDRESS     = 0x0001
	RESPONSE_ADDRESS   = 0x0002
	CHANGE_REQUEST     = 0x0003
	SOURCE_ADDRESS     = 0x0004
	CHANGED_ADDRESS    = 0x0005
	USERNAME           = 0x0006
	PASSWORD           = 0x0007
	MESSAGE_INTEGRITY  = 0x0008
	ERROR_CODE         = 0x0009
	UNKNOWN_ATTRIBUTES = 0x000a
	REFLECTED_FROM     = 0x000b
)

var (
	AttributesName = map[int]string{
		MAPPED_ADDRESS:     "MAPPED_ADDRESS",
		RESPONSE_ADDRESS:   "RESPONSE_ADDRESS",
		CHANGE_REQUEST:     "CHANGE_REQUEST",
		SOURCE_ADDRESS:     "SOURCE_ADDRESS",
		CHANGED_ADDRESS:    "CHANGED_ADDRESS",
		USERNAME:           "USERNAME",
		PASSWORD:           "PASSWORD",
		MESSAGE_INTEGRITY:  "MESSAGE_INTEGRITY",
		ERROR_CODE:         "ERROR_CODE",
		UNKNOWN_ATTRIBUTES: "UNKNOWN_ATTRIBUTES",
		REFLECTED_FROM:     "REFLECTED_FROM",
	}
)

type STUNHeaderPacket struct {
	MessageType   [2]byte
	MessageLength [2]byte
	MagicCookie   [4]byte
	TransactionID [12]byte
}

func buildRequestPacket() []byte {
	req := &STUNHeaderPacket{}

	rand.Seed(time.Now().Unix())
	tid := [12]byte{}
	for i := 0; i < 12; i++ {
		tid[i] = byte(rand.Int())
	}

	buff := make([]byte, 1024)
	binary.BigEndian.PutUint16(buff, STUNBindingRequest)
	copy(req.MessageType[:], buff)
	req.MessageLength = [2]byte{0, 0}
	binary.BigEndian.PutUint32(buff, 0x2112A442)
	copy(req.MagicCookie[:], buff)
	req.TransactionID = tid

	*(*STUNHeaderPacket)(unsafe.Pointer(&buff[0])) = *req
	size := unsafe.Sizeof(STUNHeaderPacket{})
	return buff[:size]
}

func parseResponseAttributes(buff []byte) (map[string]map[string]interface{}, error) {
	var (
		i      int
		out    = map[string]map[string]interface{}{}
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
		case MAPPED_ADDRESS, SOURCE_ADDRESS, CHANGED_ADDRESS:
			if name, ok := AttributesName[AttributeTypeValue]; ok {
				out[name] = parseAddressAttribute(buff[i : i+AttributeLength])
			}
		}
		i += AttributeLength
	}

	return out, nil
}

func parseAddressAttribute(buff []byte) map[string]interface{} {
	Reserved := util.Binary2Decimal(util.Bytes2Bits(buff[0]))
	ProtocolFamily := util.Binary2Decimal(util.Bytes2Bits(buff[1]))
	Port := util.Binary2Decimal(util.Bytes2Bits(buff[2], buff[3]))
	IP := net.IPv4(buff[4], buff[5], buff[6], buff[7]).To4().String()

	out := map[string]interface{}{
		"Reserved":       Reserved,
		"ProtocolFamily": ProtocolFamily,
		"Port":           Port,
		"IP":             IP,
	}
	return out
}
