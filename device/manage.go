package device

import "sync"

const (
	UP = iota
	DOWN
)

type Manager interface {
	SetMTU(mtu int) error
	SetAddrIPv4(addr string) error
	ChangeState(state int) error
}

type Manage struct {
	lock   sync.Mutex
	device Device
}

func NewManage(device Device) Manager {
	return &Manage{
		device: device,
	}
}
