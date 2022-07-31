package device

import (
	"os"
)

type Device interface {
	File() *os.File
	Read([]byte) (int, error)
	Write([]byte) (int, error)
	Flush() error
	MTU() (int, error)
	Name() (string, error)
	Close() error
}

type Config struct {
	IsTAP bool
}
