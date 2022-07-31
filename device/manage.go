package device

type Manager interface {
	SetMTU(mtu int) error
	SetAddrIPv4(addr string) error
}

type Manage struct {
	device Device
}

func NewManage(device Device) Manager {
	return Manage{
		device: device,
	}
}
