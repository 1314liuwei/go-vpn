package hostinfo

import (
	"fmt"

	"golang.org/x/sys/windows"
)

func init() {
	osVersionFunc = osVersionFuncWindows
}

func osVersionFuncWindows() string {
	major, minor, build := windows.RtlGetNtVersionNumbers()
	version := fmt.Sprintf("%d.%d.%d", major, minor, build)
	return version
}
