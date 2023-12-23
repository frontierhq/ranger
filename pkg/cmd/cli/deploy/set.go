package deploy

import (
	"github.com/gofrontier-com/ranger/pkg/cmd/app/deploy"
	"github.com/gofrontier-com/ranger/pkg/core"

	"github.com/spf13/cobra"
)

var (
	projectName string
	orgName     string
)

// NewCmdDeploySet creates a command to deploy a set
func NewCmdDeploySet(config *core.Config) *cobra.Command {
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
