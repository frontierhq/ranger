package deploy

import (
	"github.com/frontierdigital/ranger/core"
	"github.com/spf13/cobra"
)

// NewCmdDeploy creates a command to deploy an artifact
func NewCmdDeploy(configuration *core.Configuration) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deploy",
		Short: "Deploy an artifact",
	}

	cmd.AddCommand(NewCmdDeployManifest(configuration))

	return cmd
}
