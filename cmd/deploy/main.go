package deploy

import (
	"github.com/frontierdigital/ranger/structs"
	"github.com/spf13/cobra"
)

// NewCmdDeploy creates a command to deploy an artifact
func NewCmdDeploy(configuration *structs.Configuration) *cobra.Command {
	c := &cobra.Command{
		Use:   "deploy",
		Short: "Deploy will deploy",
		Long:  "Deploy will deploy",
	}

	c.AddCommand(NewCmdDeployManifest(configuration))

	return c
}
