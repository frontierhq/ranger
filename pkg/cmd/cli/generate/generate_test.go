package generate

import (
	"testing"

	"github.com/frontierdigital/ranger/pkg/util/config"
)

func TestNewCmdGenerate(t *testing.T) {
	config := config.Config{}
	cmd := NewCmdGenerate(&config)

	if cmd.Use != "generate" {
		t.Errorf("Use is not correct")
	}
}
