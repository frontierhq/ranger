package main

import (
	"os"

	"github.com/frontierdigital/ranger/cmd"
	"github.com/frontierdigital/ranger/core/configuration"
	"github.com/frontierdigital/ranger/core/print"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	configuration, err := configuration.LoadConfiguration()
	if err != nil {
		print.PrintlnError(err)
		os.Exit(1)
	}

	command := cmd.NewCmdRoot(configuration, version, commit, date)
	if err := command.Execute(); err != nil {
		print.PrintlnError(err)
		os.Exit(1)
	}
}
