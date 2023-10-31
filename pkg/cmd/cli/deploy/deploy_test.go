package deploy

import (
	"testing"

	"github.com/frontierdigital/ranger/pkg/core"
)

func TestNewCmdDeploy(t *testing.T) {
	config := core.Config{}
	cmd := NewCmdDeploy(&config)

	if cmd.Use != "deploy" {
		t.Errorf("Use is not correct")
	}
}
