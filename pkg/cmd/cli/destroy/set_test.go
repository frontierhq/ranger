package destroy

import (
	"testing"

	"github.com/gofrontier-com/ranger/pkg/core"
)

func TestNewCmdDestroySet(t *testing.T) {
	config := core.Config{}
	cmd := NewCmdDestroySet(&config)

	if cmd.Use != "set" {
		t.Errorf("Use is not correct")
	}
}
