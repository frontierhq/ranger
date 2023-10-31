package destroy

import (
	"testing"

	"github.com/frontierdigital/ranger/pkg/core"
)

func TestNewCmdDestroySet(t *testing.T) {
	config := core.Config{}
	cmd := NewCmdDestroySet(&config)

	if cmd.Use != "set" {
		t.Errorf("Use is not correct")
	}
}
