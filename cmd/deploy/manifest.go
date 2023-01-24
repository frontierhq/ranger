package deploy

import (
	"github.com/spf13/cobra"

	"github.com/frontierdigital/ranger/core"
	"github.com/frontierdigital/ranger/core/deploy"
)

var (
	projectName = ""
	orgName     = ""
)

// NewCmdDeployManifest creates a command to deploy a manifest
func NewCmdDeployManifest(configuration *core.Configuration) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "manifest",
		Short: "Deploy a manifest",
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := deploy.DeployManifest(configuration, projectName, orgName); err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&projectName, "project-name", "p", "", "Project name")
	cmd.Flags().StringVarP(&orgName, "organisation-name", "o", "", "Organisation name")

	return cmd
}
