package device

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateTUN(t *testing.T) {
	name := "testTun"
	tun, err := CreateTUN(name)
	if err != nil {
		return
	}
	defer tun.Close()

	byName, err := net.InterfaceByName(name)
	assert.NotNil(t, byName)
}
