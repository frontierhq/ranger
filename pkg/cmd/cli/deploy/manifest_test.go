package deploy

import (
	"testing"

	"github.com/frontierdigital/ranger/pkg/cmd/app/type/config"
)

func TestNewCmdDeployManifest(t *testing.T) {
	config := config.Config{}
	cmd := NewCmdDeployManifest(&config)

	if cmd.Use != "manifest" {
		t.Errorf("Use is not correct")
	}
}
