package device

import (
	"io"
	"sync"
)

const (
	UP = iota
	DOWN
)

type Manager interface {
	io.ReadWriteCloser
	SetMTU(mtu int) error
	SetAddrIPv4(addr string) error
	ChangeState(state int) error
}

type Manage struct {
	lock   sync.Mutex
	device Device
}

func (m *Manage) Read(buff []byte) (int, error) {
	return m.device.Read(buff)
}

func (m *Manage) Write(buff []byte) (int, error) {
	return m.device.Write(buff)
}

func (m *Manage) Close() error {
	return m.device.Close()
}

func NewManage(device Device) Manager {
	return &Manage{
		device: device,
	}
}
