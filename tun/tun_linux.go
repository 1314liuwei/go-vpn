package tun

import (
	"bytes"
	"fmt"
	"os"
	"syscall"
	"unsafe"

	"golang.org/x/sys/unix"
)

const (
	defaultCloneTUNPath = "/dev/net/tun"
)

// ifreq 占用 40个字节: https://man7.org/linux/man-pages/man7/netdevice.7.html
type ifReq [40]byte

// ioctl 方法：https://man7.org/linux/man-pages/man2/ioctl.2.html
func ioctl(fd uintptr, request uintptr, argp uintptr) error {
	_, _, err := unix.Syscall(unix.SYS_IOCTL, fd, request, argp)
	if err != 0 {
		return os.NewSyscallError("ioctl: ", err)
	}
	return nil
}

var _ Device = new(NativeTUN)

type NativeTUN struct {
	tunFile     *os.File
	name        string
	index       int32
	netlinkSock int
}

func (tun NativeTUN) File() *os.File {
	//TODO implement me
	panic("implement me")
}

func (tun NativeTUN) Read(bytes []byte) (int, error) {
	//TODO implement me
	panic("implement me")
}

func (tun NativeTUN) Write(bytes []byte) (int, error) {
	//TODO implement me
	panic("implement me")
}

func (tun NativeTUN) Flush() error {
	//TODO implement me
	panic("implement me")
}

func (tun NativeTUN) MTU() (int, error) {
	//TODO implement me
	panic("implement me")
}

func (tun NativeTUN) Name() (string, error) {
	if tun.name != "" {
		return tun.name, nil
	} else {
		return tun.getNameFromInterface()
	}
}

func (tun NativeTUN) Close() error {
	err := tun.tunFile.Close()
	if err != nil {
		return err
	}
	return nil
}

func (tun NativeTUN) setMTU(mtu int) error {
	name, err := tun.Name()
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

func (tun NativeTUN) getNameFromInterface() (string, error) {
	conn, err := tun.tunFile.SyscallConn()
	if err != nil {
		return "", err
	}

	var ifr ifReq
	var errno syscall.Errno
	// 获取 TUN 设备信息
	err = conn.Control(func(fd uintptr) {
		_, _, errno = unix.Syscall(
			unix.SYS_IOCTL,
			fd,
			uintptr(unix.TUNGETIFF),
			uintptr(unsafe.Pointer(&ifr)),
		)
	})
	if err != nil || errno != 0 {
		return "", fmt.Errorf("failed to get name of TUN device: %w", err)
	}

	name := ifr[:]
	if i := bytes.IndexByte(name, 0); i != -1 {
		name = name[:i]
	}

	tun.name = string(name[:])
	return tun.name, nil
}

func CreateTUN(name string, mtu int) (Device, error) {
	fd, err := unix.Open(defaultCloneTUNPath, os.O_RDWR, 0)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("CreateTUN(%q) failed; %s does not exist", name, defaultCloneTUNPath)
		}
		return nil, err
	}

	var ifr ifReq
	nameBytes := []byte(name)
	// 网卡名称不能过长
	if len(nameBytes) >= unix.IFNAMSIZ {
		if err := unix.Close(fd); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("interface name too long: %w", unix.ENAMETOOLONG)
	}

	// 拷贝名字
	copy(ifr[:], nameBytes)
	// 拷贝 flags
	*(*uint16)(unsafe.Pointer(&ifr[unix.IFNAMSIZ])) = unix.IFF_TUN
	err = ioctl(uintptr(fd), unix.TUNSETIFF, uintptr(unsafe.Pointer(&ifr)))
	if err != nil {
		return nil, err
	}

	err = unix.SetNonblock(fd, true)
	if err != nil {
		if err := unix.Close(fd); err != nil {
			return nil, err
		}
		return nil, err
	}

	file := os.NewFile(uintptr(fd), defaultCloneTUNPath)
	return CreateTUNFromFile(file, mtu)
}

func CreateTUNFromFile(file *os.File, mtu int) (Device, error) {
	tun := &NativeTUN{
		tunFile: file,
	}

	name, err := tun.Name()
	if err != nil {
		return nil, err
	}

	// 设置网卡设备序列号
	tun.index, err = getIFIndex(name)
	if err != nil {
		return nil, err
	}

	tun.netlinkSock, err = createNetlinkSocket()
	if err != nil {
		unix.Close(tun.netlinkSock)
		return nil, err
	}

	err = tun.setMTU(mtu)
	if err != nil {
		return nil, err
	}

	return tun, nil
}

// getIFIndex 获取设备序列号
func getIFIndex(name string) (int32, error) {
	fd, err := unix.Socket(unix.AF_INET, unix.SOCK_DGRAM, 0)
	if err != nil {
		return 0, err
	}
	defer unix.Close(fd)

	var ifr ifReq
	copy(ifr[:], name)
	// 获取设备序列号
	err = ioctl(uintptr(fd), unix.SIOCGIFINDEX, uintptr(unsafe.Pointer(&ifr[0])))
	if err != nil {
		return 0, err
	}
	// 获取序列号，参考 ifreq 结构体
	return *(*int32)(unsafe.Pointer(&ifr[unix.IFNAMSIZ])), nil
}

func createNetlinkSocket() (int, error) {
	sock, err := unix.Socket(unix.AF_NETLINK, unix.SOCK_RAW, unix.NETLINK_ROUTE)
	if err != nil {
		return -1, err
	}

	saddr := &unix.SockaddrNetlink{
		Family: unix.AF_NETLINK,
		Groups: unix.RTMGRP_LINK | unix.RTMGRP_IPV4_IFADDR | unix.RTMGRP_IPV6_IFADDR,
	}
	err = unix.Bind(sock, saddr)
	if err != nil {
		return -1, err
	}

	return sock, nil
}
