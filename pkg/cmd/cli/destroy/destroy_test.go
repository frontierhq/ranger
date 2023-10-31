package destroy

import (
	"testing"

	"github.com/frontierdigital/ranger/pkg/core"
)

func TestNewCmdDestroy(t *testing.T) {
	config := core.Config{}
	cmd := NewCmdDestroy(&config)

	if cmd.Use != "destroy" {
		t.Errorf("Use is not correct")
	}
}
