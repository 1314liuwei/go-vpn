package protocol

type ICMP struct {
	Type             byte
	Code             byte
	CheckSum         [2]byte
	IdentifierBE     byte
	IdentifierLE     byte
	SequenceNumberBE byte
	SequenceNumberLE byte
	Timestamp        [8]byte
	Data             []byte
}

// 0 0 8 0 69 0 0 84 124 194 0 0 64 1 122 147 192 168 1 2 192 168 1 1 0 0 177 117 154 249 0 11 96 10 18 99 0 0 0 0 126 69 4 0 0 0 0 0 16 17 18 19 20 21 22 23 24 25 26 27 28 29 30 31 32 33 34 35 36 37 38 39 40 41 42 43 44 45 46 47 48 49 50 51 52 53 54 55
func icmpParse(packet []byte) {

}
