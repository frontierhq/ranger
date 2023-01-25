package deploy

import (
	"testing"

	"github.com/frontierdigital/ranger/pkg/cmd/app"
)

func TestNewCmdDeployManifest(t *testing.T) {
	config := app.Config{}
	cmd := NewCmdDeployManifest(&config)

	if cmd.Use != "manifest" {
		t.Errorf("Use is not correct")
	}
}
