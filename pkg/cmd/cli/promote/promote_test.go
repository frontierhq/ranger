package promote

import (
	"testing"

	"github.com/gofrontier-com/ranger/pkg/core"
)

func TestNewCmdPromote(t *testing.T) {
	config := core.Config{}
	cmd := NewCmdPromote(&config)

	if cmd.Use != "promote" {
		t.Errorf("Use is not correct")
	}
}
