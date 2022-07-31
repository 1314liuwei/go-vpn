package device

import (
	"golang.org/x/sys/unix"
	"net"
	"unsafe"
)

func (m *Manage) SetMTU(mtu int) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	name, err := m.device.Name()
	if err != nil {
		return err
	}

	// 获取 socket 文件句柄
	fd, err := unix.Socket(unix.AF_INET, unix.SOCK_DGRAM, 0)
	if err != nil {
		return err
	}
	defer unix.Close(fd)

	// 设置 mtu
	var ifr ifReq
	copy(ifr[:], name)
	*(*uint32)(unsafe.Pointer(&ifr[unix.IFNAMSIZ])) = uint32(mtu)
	err = ioctl(uintptr(fd), unix.SIOCSIFMTU, uintptr(unsafe.Pointer(&ifr)))
	if err != nil {
		return err
	}

	return nil
}

func (m *Manage) SetAddrIPv4(addr string) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	name, err := m.device.Name()
	if err != nil {
		return err
	}

	fd, err := unix.Socket(unix.AF_INET, unix.SOCK_STREAM, 0)
	if err != nil {
		return err
	}
	defer unix.Close(fd)

	// 配置 ip 地址
	var ifr ifReq
	var address [4]byte

	copy(ifr[:], name)
	copy(address[:], net.ParseIP(addr))
	sock := unix.RawSockaddrInet4{
		Family: unix.AF_INET,
		Addr:   address,
	}
	*(*unix.RawSockaddrInet4)(unsafe.Pointer(&ifr[unix.IFNAMSIZ])) = sock

	err = ioctl(uintptr(fd), unix.SIOCSIFADDR, uintptr(unsafe.Pointer(&ifr)))
	if err != nil {
		return err
	}

	return nil
}

func (m *Manage) ChangeState(state int) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	var ifr ifReq
	name, err := m.device.Name()
	if err != nil {
		return err
	}
	copy(ifr[:], name)
	fd, err := unix.Socket(unix.AF_INET, unix.SOCK_STREAM, 0)
	if err != nil {
		return err
	}
	defer unix.Close(fd)

	switch state {
	case UP:
		// 启动网卡
		*(*uint16)(unsafe.Pointer(&ifr[unix.IFNAMSIZ])) = unix.IFF_UP | unix.IFF_RUNNING
	case DOWN:
		// 关闭网卡
		*(*uint16)(unsafe.Pointer(&ifr[unix.IFNAMSIZ])) = ^uint16(unix.IFF_UP)
	}

	err = ioctl(uintptr(fd), unix.SIOCSIFFLAGS, uintptr(unsafe.Pointer(&ifr)))
	if err != nil {
		return err
	}

	err = ioctl(uintptr(fd), unix.SIOCSIFFLAGS, uintptr(unsafe.Pointer(&ifr)))
	if err != nil {
		return err
	}
	return nil
}
