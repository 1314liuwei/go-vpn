package hostinfo

import (
	"bufio"
	"os"
	"runtime"
	"strings"
)

type HostInfo struct {
	Hostname       string
	OS             string
	OSVersion      string
	LinuxDesktopOS bool
}

type probe struct{}

func New() HostInfo {
	return HostInfo{}
}

var (
	osVersionFunc func() string
)

func (p probe) hostName() string {
	hostname, err := os.Hostname()
	if err != nil {
		return ""
	}
	return hostname
}

func (p probe) os() string {
	switch runtime.GOOS {
	case "ios":
		return "iOS"
	case "darwin":
		return "macOS"
	default:
		return runtime.GOOS
	}
}

func (p probe) osVersion() string {
	if osVersionFunc != nil {
		return osVersionFunc()
	}
	return ""
}

func (p probe) linuxDesktopOS() bool {
	if runtime.GOOS != "linux" {
		return false
	}

	procFile := "/proc/net/unix"
	file, err := os.Open(procFile)
	if err != nil {
		return false
	}

	desktopProc := []string{
		" @/tmp/dbus-",
		".X11-unix",
		"/wayland-1",
	}
	buff := bufio.NewScanner(file)
	for {
		if !buff.Scan() {
			break
		}
		line := buff.Text()
		for i := 0; i < len(desktopProc); i++ {
			if strings.Contains(line, desktopProc[i]) {
				return true
			}
		}
	}

	return false
}
