package cmd

import (
	"os"

	"github.com/spf13/cobra"

	vers "github.com/frontierdigital/ranger/cmd/version"
	"github.com/frontierdigital/ranger/core/configuration"
)

func NewCmdRoot(configuration *configuration.Configuration, version string, commit string, date string) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:                   "ranger",
		DisableFlagsInUseLine: true,
		Short:                 "ranger is the command line tool for Ranger",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				cmd.Help()
				os.Exit(0)
			}

			return nil
		},
	}

	rootCmd.AddCommand(vers.NewCmdVersion(version, commit, date))

	return rootCmd
}
