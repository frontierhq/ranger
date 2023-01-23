package cmd

import (
	"github.com/spf13/cobra"

	"github.com/frontierdigital/ranger/cmd/deploy"
	vers "github.com/frontierdigital/ranger/cmd/version"
	"github.com/frontierdigital/ranger/structs"
)

func NewCmdRoot(configuration *structs.Configuration, version string, commit string, date string) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:                   "ranger",
		DisableFlagsInUseLine: true,
		Short:                 "ranger is the command line tool for Ranger",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmd.Help(); err != nil {
				return err
			}

			return nil
		},
	}

	rootCmd.AddCommand(deploy.NewCmdDeploy(configuration))
	rootCmd.AddCommand(vers.NewCmdVersion(version, commit, date))

	return rootCmd
}
