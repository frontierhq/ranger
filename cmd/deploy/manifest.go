package deploy

import (
	"github.com/spf13/cobra"

	"github.com/frontierdigital/ranger/core/deploy"
	"github.com/frontierdigital/ranger/structs"
)

var (
	projectName = ""
	orgName     = ""
)

// NewCmdDeployManifest creates a command to deploy a manifest
func NewCmdDeployManifest(configuration *structs.Configuration) *cobra.Command {
	c := &cobra.Command{
		Use:   "manifest",
		Short: "deploy a manifest",
		Long:  "deploy a manifest",
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := deploy.DeployManifest(configuration, projectName, orgName); err != nil {
				return err
			}

			return nil
		},
	}

	c.Flags().StringVarP(&projectName, "project-name", "p", "", "Project name")
	c.Flags().StringVarP(&orgName, "organisation-name", "o", "", "Organisation name")

	return c
}
