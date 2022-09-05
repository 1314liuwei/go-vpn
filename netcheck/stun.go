package netcheck

import (
	"encoding/binary"
	"math/rand"
	"time"
	"unsafe"
)

const (
	STUNBindingRequest       = 0x0001
	STUNBindingResponse      = 0x0101
	STUNBindingErrorResponse = 0x0111
)

type STUNHeaderPacket struct {
	MessageType   [2]byte
	MessageLength [2]byte
	MagicCookie   [4]byte
	TransactionID [12]byte
}

type ResponsePacket struct {
	STUNHeaderPacket
	Attributes map[string]map[string]string
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

func parseResponsePacket(buff []byte) {

}
