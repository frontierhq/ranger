package generate

import (
	"testing"

	"github.com/frontierdigital/ranger/pkg/util/config"
)

func TestNewCmdGenerateDocs(t *testing.T) {
	config := config.Config{}
	cmd := NewCmdGenerateDocs(&config)

	if cmd.Use != "generate" {
		t.Errorf("Use is not correct")
	}
}
