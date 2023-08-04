package deploy

import (
	"testing"

	"github.com/frontierdigital/ranger/pkg/util/config"
)

func TestNewCmdDeploySet(t *testing.T) {
	config := config.Config{}
	cmd := NewCmdDeploySet(&config)

	if cmd.Use != "set" {
		t.Errorf("Use is not correct")
	}
}
