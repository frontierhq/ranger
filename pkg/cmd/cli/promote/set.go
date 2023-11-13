package promote

import (
	"github.com/frontierdigital/ranger/pkg/cmd/app/promote"
	"github.com/frontierdigital/ranger/pkg/core"

	"github.com/spf13/cobra"
)

var (
	projectName     = ""
	nextEnvironment = ""
	orgName         = ""
)

// NewCmdPromoteSet creates a command to promote a set
func NewCmdPromoteSet(config *core.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set",
		Short: "Promote a set",
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := promote.PromoteSet(config, projectName, orgName); err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&projectName, "project-name", "p", "", "Project name")
	cmd.Flags().StringVarP(&orgName, "organisation-name", "o", "", "Organisation name")

	return cmd
}
