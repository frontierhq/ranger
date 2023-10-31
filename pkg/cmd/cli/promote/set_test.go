package promote

import (
	"testing"

	"github.com/frontierdigital/ranger/pkg/core"
)

func TestNewCmdPromoteSet(t *testing.T) {
	config := core.Config{}
	cmd := NewCmdPromoteSet(&config)

	if cmd.Use != "set" {
		t.Errorf("Use is not correct")
	}
}
