package tun

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

func New(conf Config) (Device, error) {

	if conf.IsTAP {

	} else {
	}

	return nil, nil
}
