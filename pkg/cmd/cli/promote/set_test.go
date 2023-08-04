package promote

import (
	"testing"

	"github.com/frontierdigital/ranger/pkg/util/config"
)

func TestNewCmdPromoteSet(t *testing.T) {
	config := config.Config{}
	cmd := NewCmdPromoteSet(&config)

	if cmd.Use != "set" {
		t.Errorf("Use is not correct")
	}
}
