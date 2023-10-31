package generate

import (
	"testing"

	"github.com/frontierdigital/ranger/pkg/core"
)

func TestNewCmdGenerateDocs(t *testing.T) {
	config := core.Config{}
	cmd := NewCmdGenerateDocs(&config)

	if cmd.Use != "docs" {
		t.Errorf("Use is not correct")
	}
}
