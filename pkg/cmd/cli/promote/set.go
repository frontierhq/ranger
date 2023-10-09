package promote

import (
	"github.com/frontierdigital/ranger/pkg/cmd/app/promote"
	"github.com/frontierdigital/ranger/pkg/util/config"

	"github.com/spf13/cobra"
)

var (
	projectName     = ""
	nextEnvironment = ""
	orgName         = ""
)

// NewCmdPromoteSet creates a command to promote a set
func NewCmdPromoteSet(config *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set",
		Short: "Promote a set",
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := promote.PromoteSet(config, projectName, orgName, nextEnvironment); err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&projectName, "project-name", "p", "", "Project name")
	cmd.Flags().StringVarP(&nextEnvironment, "next-environment", "n", "", "Next environment")
	cmd.Flags().StringVarP(&orgName, "organisation-name", "o", "", "Organisation name")

	return cmd
}
