package device

import "testing"

func TestWrap_SetMTU(t *testing.T) {
}

func TestWrap_SetAddrIPv4(t *testing.T) {
	name := "tun0"
	tun, err := CreateTUN(name)
	if err != nil {
		return
	}

	wrap := NewManage(tun)
	err = wrap.SetAddrIPv4("192.168.111.1")
	if err != nil {
		return
	}
}
