package ranger

import (
	"os"

	"github.com/frontierdigital/ranger/pkg/cmd/cli/deploy"
	vers "github.com/frontierdigital/ranger/pkg/cmd/cli/version"
	"github.com/frontierdigital/ranger/pkg/util/config"
	"github.com/frontierdigital/utils/output"

	"github.com/spf13/cobra"
)

func NewRootCmd(version string, commit string, date string) *cobra.Command {
	config, err := config.LoadConfig()
	if err != nil {
		output.PrintlnError(err)
		os.Exit(1)
	}

	cmd := &cobra.Command{
		Use:                   "ranger",
		DisableFlagsInUseLine: true,
		Short:                 "ranger is the command line tool for Ranger",
	}

	cmd.AddCommand(deploy.NewCmdDeploy(config))
	cmd.AddCommand(vers.NewCmdVersion(version, commit, date))

	return cmd
}
