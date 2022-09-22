package hostinfo

import (
	"os"
	"runtime"
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
	return true
}
