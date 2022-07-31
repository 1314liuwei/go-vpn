package device

import (
	"bytes"
	"fmt"
	"os"
	"syscall"
	"unsafe"

	"golang.org/x/sys/unix"
)

const (
	defaultCloneTUNPath = "/dev/net/device"
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
	return tun.tunFile
}

func (tun NativeTUN) Read(bytes []byte) (int, error) {
	return tun.tunFile.Read(bytes)
}

func (tun NativeTUN) Write(bytes []byte) (int, error) {
	return tun.tunFile.Write(bytes)
}

func (tun NativeTUN) Flush() error {
	return nil
}

func (tun NativeTUN) MTU() (int, error) {
	name, err := tun.Name()
	if err != nil {
		return -1, err
	}

	fd, err := unix.Socket(unix.AF_INET, unix.SOCK_DGRAM, 0)
	if err != nil {
		return 0, err
	}
	defer unix.Close(fd)

	var ifr ifReq
	copy(ifr[:], name)
	err = ioctl(uintptr(fd), unix.SIOCGIFMTU, uintptr(unsafe.Pointer(&ifr)))
	if err != nil {
		return -1, err
	}
	return int(*(*int32)(unsafe.Pointer(&ifr[unix.IFNAMSIZ]))), nil
}

func (tun NativeTUN) Name() (string, error) {
	if tun.name != "" {
		return tun.name, nil
	} else {
		return tun.getNameFromSystem()
	}
}

func (tun NativeTUN) Close() error {
	err := tun.tunFile.Close()
	if err != nil {
		return err
	}
	return nil
}

func (tun NativeTUN) getNameFromSystem() (string, error) {
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

func CreateTUN(name string) (Device, error) {
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
	// 设置 flags
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

	// 将文件描述符转换为 *os.File 对象
	file := os.NewFile(uintptr(fd), defaultCloneTUNPath)
	return CreateTUNFromFile(file)
}

func CreateTUNFromFile(file *os.File) (Device, error) {
	tun := &NativeTUN{
		tunFile: file,
	}

	// 获取网卡名字
	name, err := tun.Name()
	if err != nil {
		return nil, err
	}

	// 获取网卡唯一索引号
	tun.index, err = getIFIndex(name)
	if err != nil {
		return nil, err
	}

	tun.netlinkSock, err = createNetlinkSocket()
	if err != nil {
		unix.Close(tun.netlinkSock)
		return nil, err
	}

	return tun, nil
}

// getIFIndex 获取网卡唯一索引号
func getIFIndex(name string) (int32, error) {
	fd, err := unix.Socket(unix.AF_INET, unix.SOCK_DGRAM, 0)
	if err != nil {
		return 0, err
	}
	defer unix.Close(fd)

	var ifr ifReq
	copy(ifr[:], name)
	// 获取网卡唯一索引号
	err = ioctl(uintptr(fd), unix.SIOCGIFINDEX, uintptr(unsafe.Pointer(&ifr[0])))
	if err != nil {
		return 0, err
	}
	// 解析得到索引号，参考 ifreq 结构体
	return *(*int32)(unsafe.Pointer(&ifr[unix.IFNAMSIZ])), nil
}

func createNetlinkSocket() (int, error) {
	// 创建一个 NETLINK_ROUTE 类型的 Netlink 套接字
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
