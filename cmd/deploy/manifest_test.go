package deploy

import (
	"testing"
)

func TestNewCmdDeployManifest(t *testing.T) {
	cmd := NewCmdDeployManifest()

	if cmd.Use != "manifest" {
		t.Errorf("Use is not correct")
	}
}
