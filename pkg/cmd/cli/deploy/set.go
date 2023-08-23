package deploy

import (
	"github.com/frontierdigital/ranger/pkg/cmd/app/deploy"
	"github.com/frontierdigital/ranger/pkg/util/config"

	"github.com/spf13/cobra"
)

var (
	projectName = ""
	orgName     = ""
)

// NewCmdDeploySet creates a command to deploy a set
func NewCmdDeploySet(config *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set",
		Short: "Deploy a set",
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := deploy.DeploySet(config, projectName, orgName); err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&projectName, "project-name", "p", "", "Project name")
	cmd.Flags().StringVarP(&orgName, "organisation-name", "o", "", "Organisation name")

	return cmd
}
