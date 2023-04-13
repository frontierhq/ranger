package main

import (
	"os"

	"github.com/frontierdigital/ranger/pkg/cmd/ranger"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	command := ranger.NewRootCmd(version, commit, date)
	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
