package cmd

import (
	"testing"

	"github.com/frontierdigital/ranger/core/structs"
)

func TestNewCmdRoot(t *testing.T) {
	configuration := structs.Configuration{}
	cmd := NewCmdRoot(&configuration, "0.0.0", "commitid", "date")

	if cmd.Use != "ranger" {
		t.Errorf("Use is not correct")
	}
}
