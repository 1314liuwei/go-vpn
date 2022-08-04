package config

import (
	"testing"

	"github.com/gogf/gf/v2/frame/g"
)

func TestNew(t *testing.T) {
	config := New()
	g.Dump(config)
}
