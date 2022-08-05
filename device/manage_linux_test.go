package device

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWrap_SetMTU(t *testing.T) {
}

func TestManage_ChangeState(t *testing.T) {
	name := "tun0"
	tun, err := CreateTUN(name)
	defer tun.Close()
	assert.Nil(t, err)

	wrap := NewManage(tun)
	err = wrap.ChangeState(UP)
	assert.Nil(t, err)

	fmt.Println("up: ")
	time.Sleep(60 * time.Second)
}

func TestWrap_SetAddrIPv4(t *testing.T) {
	name := "tun0"
	tun, err := CreateTUN(name)
	defer tun.Close()
	assert.Nil(t, err)

	wrap := NewManage(tun)
	err = wrap.ChangeState(UP)
	assert.Nil(t, err)

	err = wrap.SetAddrIPv4("192.168.111.1")
	assert.Nil(t, err)

	fmt.Println("set addr: ")
	time.Sleep(60 * time.Second)
}

func TestTUNPing(t *testing.T) {
	name := "tun0"
	tun, err := CreateTUN(name)
	defer tun.Close()
	assert.Nil(t, err)

	wrap := NewManage(tun)
	wrap.SetMTU(1280)
	err = wrap.ChangeState(UP)
	assert.Nil(t, err)

	err = wrap.SetAddrIPv4("192.168.111.1/24")
	if err != nil {
		panic(err)
	}
	time.Sleep(10 * time.Second)

	for {
		var (
			buf [1500]byte
			ip  [4]byte
		)

		read, err := tun.Read(buf[:])
		if err != nil {
			panic(err)
		}
		fmt.Println(buf[:read+1])

		copy(ip[:], buf[16:20])
		copy(buf[16:20], buf[20:24])
		copy(buf[20:24], ip[:])
		buf[24] = 0
		buf[26] += 8

		write, err := tun.Write(buf[:read])
		if err != nil {
			panic(err)
		}
		fmt.Println(write)
	}

}
