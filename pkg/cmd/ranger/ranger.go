package ranger

import (
	"os"

	"github.com/gofrontier-com/go-utils/output"
	"github.com/gofrontier-com/ranger/pkg/cmd/cli/deploy"
	"github.com/gofrontier-com/ranger/pkg/cmd/cli/destroy"
	"github.com/gofrontier-com/ranger/pkg/cmd/cli/generate"
	"github.com/gofrontier-com/ranger/pkg/cmd/cli/promote"
	vers "github.com/gofrontier-com/ranger/pkg/cmd/cli/version"
	"github.com/gofrontier-com/ranger/pkg/util/config"

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
	cmd.AddCommand(destroy.NewCmdDestroy(config))
	cmd.AddCommand(generate.NewCmdGenerate(config))
	cmd.AddCommand(promote.NewCmdPromote(config))
	cmd.AddCommand(vers.NewCmdVersion(version, commit, date))

	return cmd
}
