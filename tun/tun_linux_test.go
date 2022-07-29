package tun

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateTUN(t *testing.T) {
	name := "testTun"
	tun, err := CreateTUN(name, 1280)
	if err != nil {
		return
	}
	defer tun.Close()

	byName, err := net.InterfaceByName(name)
	assert.NotNil(t, byName)
	assert.Equal(t, byName.MTU, 1280)
}
