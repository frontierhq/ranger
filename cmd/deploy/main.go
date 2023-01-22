package deploy

import (
	"github.com/spf13/cobra"

	"github.com/frontierdigital/ranger/core/configuration"
)

// NewCmdDeploy creates a new deploy command
func NewCmdDeploy(cfg *configuration.Configuration) *cobra.Command {
	c := &cobra.Command{
		Use:   "deploy",
		Short: "Deploy will deploy",
		Long:  "Deploy will deploy",
		RunE: func(cmd *cobra.Command, _ []string) error {
			if err := cmd.Help(); err != nil {
				return err
			}

			return nil
		},
	}

	c.AddCommand(NewCmdDeployManifest(cfg))

	return c
}
