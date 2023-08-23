package promote

import (
	"testing"

	"github.com/frontierdigital/ranger/pkg/util/config"
)

func TestNewCmdPromote(t *testing.T) {
	config := config.Config{}
	cmd := NewCmdPromote(&config)

	if cmd.Use != "promote" {
		t.Errorf("Use is not correct")
	}
}
