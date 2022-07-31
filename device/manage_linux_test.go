package device

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
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
