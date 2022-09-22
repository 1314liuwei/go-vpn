package hostinfo

import "io/ioutil"

func init() {
	osVersionFunc = osVersionFuncLinux
}

func osVersionFuncLinux() string {
	propFile := "/etc/os-release"
	file, err := ioutil.ReadFile(propFile)
	if err != nil {
		return ""
	}
	return string(file)
}
