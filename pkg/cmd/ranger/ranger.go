package ranger

import (
	"os"

	vers "github.com/frontierdigital/ranger/pkg/cmd/cli/version"
	"github.com/frontierdigital/ranger/pkg/util/configuration"
	"github.com/frontierdigital/ranger/pkg/util/output"
	"github.com/spf13/cobra"
)

func NewRootCmd(version string, commit string, date string) *cobra.Command {
	_, err := configuration.LoadConfiguration()
	if err != nil {
		output.PrintlnError(err)
		os.Exit(1)
	}

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

	rootCmd.AddCommand(vers.NewCmdVersion(version, commit, date))

	return rootCmd
}
