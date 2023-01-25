package main

import (
	"os"

	"github.com/frontierdigital/ranger/pkg/cmd/ranger"
	"github.com/frontierdigital/ranger/pkg/util/output"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	command := ranger.NewRootCmd(version, commit, date)
	if err := command.Execute(); err != nil {
		output.PrintlnError(err)
		os.Exit(1)
	}
}
