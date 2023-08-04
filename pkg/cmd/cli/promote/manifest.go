package promote

import (
	"github.com/frontierdigital/ranger/pkg/cmd/app/promote"
	"github.com/frontierdigital/ranger/pkg/util/config"

	"github.com/spf13/cobra"
)

var (
	projectName       = ""
	targetEnvironment = ""
	orgName           = ""
)

// NewCmdPromoteManifest creates a command to promote a manifest
func NewCmdPromoteManifest(config *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "manifest",
		Short: "Promote a manifest",
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := promote.PromoteManifest(config, projectName, orgName, targetEnvironment); err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&projectName, "project-name", "p", "", "Project name")
	cmd.Flags().StringVarP(&targetEnvironment, "target-environment", "n", "", "Target environment")
	cmd.Flags().StringVarP(&orgName, "organisation-name", "o", "", "Organisation name")

	return cmd
}
