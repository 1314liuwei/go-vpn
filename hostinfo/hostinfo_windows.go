package hostinfo

import (
	"fmt"

	"golang.org/x/sys/windows"
)

func init() {
	osVersionFunc = osVersionFuncWindows
}

func osVersionFuncWindows() string {
	version, err := windows.GetVersion()
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%d.%d (%d)", byte(version), uint8(version>>8), version>>16)
}
