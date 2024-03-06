package generate

import (
	"testing"

	"github.com/gofrontier-com/ranger/pkg/core"
)

func TestNewCmdGenerate(t *testing.T) {
	config := core.Config{}
	cmd := NewCmdGenerate(&config)

	if cmd.Use != "generate" {
		t.Errorf("Use is not correct")
	}
}
