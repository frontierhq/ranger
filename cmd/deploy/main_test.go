package deploy

import (
	"testing"

	"github.com/frontierdigital/ranger/core"
)

func TestNewCmdDeploy(t *testing.T) {
	configuration := core.Configuration{}
	cmd := NewCmdDeploy(&configuration)

	if cmd.Use != "deploy" {
		t.Errorf("Use is not correct")
	}
}
