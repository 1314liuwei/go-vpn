package hostinfo

import (
	"log"
	"testing"
)

func TestNew(t *testing.T) {
	log.Println(osVersionFuncWindows())
}
