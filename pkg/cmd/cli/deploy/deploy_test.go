package deploy

import (
	"testing"

	"github.com/frontierdigital/ranger/pkg/cmd/app"
)

func TestNewCmdDeploy(t *testing.T) {
	config := app.Config{}
	cmd := NewCmdDeploy(&config)

	if cmd.Use != "deploy" {
		t.Errorf("Use is not correct")
	}
}
