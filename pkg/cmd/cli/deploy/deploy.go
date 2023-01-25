package deploy

import (
	"github.com/frontierdigital/ranger/pkg/cmd/app"
	"github.com/spf13/cobra"
)

// NewCmdDeploy creates a command to deploy an artifact
func NewCmdDeploy(config *app.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deploy",
		Short: "Deploy an artifact",
	}

	cmd.AddCommand(NewCmdDeployManifest(config))

	return cmd
}
