package deploy

import (
	"testing"

	"github.com/gofrontier-com/ranger/pkg/core"
)

func TestNewCmdDeploySet(t *testing.T) {
	config := core.Config{}
	cmd := NewCmdDeploySet(&config)

	if cmd.Use != "set" {
		t.Errorf("Use is not correct")
	}
}
