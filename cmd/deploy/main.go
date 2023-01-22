package deploy

import (
	"github.com/spf13/cobra"
)

// NewCmdDeploy creates a new deploy command
func NewCmdDeploy() *cobra.Command {
	c := &cobra.Command{
		Use:   "deploy",
		Short: "Deploy will deploy",
		Long:  "Deploy will deploy",
		RunE: func(cmd *cobra.Command, _ []string) error {
			if err := cmd.Help(); err != nil {
				return err
			}

			return nil

			return nil
		},
	}

	c.AddCommand(NewCmdDeployManifest())

	return c
}
