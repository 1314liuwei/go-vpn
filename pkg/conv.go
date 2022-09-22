package pkg

func Bytes2Bits(data ...byte) []byte {
	dst := make([]byte, 0)
	for _, datum := range data {
		for i := 0; i < 8; i++ {
			move := uint(7 - i)
			dst = append(dst, (datum>>move)&1)
		}
	}
	return dst
}

func Binary2Decimal(buff []byte) int {
	out := 0
	for i := len(buff) - 1; i >= 0; i-- {
		// 扩大类型容量
		num := int(buff[i])
		out += num << (len(buff) - 1 - i)
	}
	return out
}
