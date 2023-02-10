package deploy

import (
	"testing"

	"github.com/frontierdigital/ranger/pkg/util/config"
)

func TestNewCmdDeploy(t *testing.T) {
	config := config.Config{}
	cmd := NewCmdDeploy(&config)

	if cmd.Use != "deploy" {
		t.Errorf("Use is not correct")
	}
}
