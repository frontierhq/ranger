package deploy

import (
	"testing"
)

func TestNewCmdDeploy(t *testing.T) {
	cmd := NewCmdDeploy()

	if cmd.Use != "deploy" {
		t.Errorf("Use is not correct")
	}
}
