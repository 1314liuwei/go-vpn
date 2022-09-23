package hostinfo

import (
	"io/ioutil"
	"regexp"
)

func init() {
	osVersionFunc = osVersionFuncLinux
}

func osVersionFuncLinux() string {
	propFile := "/etc/os-release"
	file, err := ioutil.ReadFile(propFile)
	if err != nil {
		return ""
	}

	compile, err := regexp.Compile(`PRETTY_NAME="(.*?)"`)
	if err != nil {
		return ""
	}
	result := compile.FindSubmatch(file)
	if len(result) > 1 {
		return string(result[1])
	}
	return ""
}
