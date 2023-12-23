package deploy

import (
	"github.com/gofrontier-com/ranger/pkg/core"
	"github.com/spf13/cobra"
)

// NewCmdDeploy creates a command to deploy an artifact
func NewCmdDeploy(config *core.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deploy",
		Short: "Deploy an artifact",
	}

	cmd.AddCommand(NewCmdDeploySet(config))

	return cmd
}
