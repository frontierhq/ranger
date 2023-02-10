package deploy

import (
	"github.com/frontierdigital/ranger/pkg/util/config"

	"github.com/spf13/cobra"
)

// NewCmdDeploy creates a command to deploy an artifact
func NewCmdDeploy(config *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deploy",
		Short: "Deploy an artifact",
	}

	cmd.AddCommand(NewCmdDeployManifest(config))

	return cmd
}
