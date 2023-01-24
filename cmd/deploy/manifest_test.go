package deploy

import (
	"testing"

	"github.com/frontierdigital/ranger/core"
)

func TestNewCmdDeployManifest(t *testing.T) {
	configuration := core.Configuration{}
	cmd := NewCmdDeployManifest(&configuration)

	if cmd.Use != "manifest" {
		t.Errorf("Use is not correct")
	}
}
