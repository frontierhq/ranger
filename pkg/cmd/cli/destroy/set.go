package destroy

import (
	"github.com/gofrontier-com/ranger/pkg/cmd/app/destroy"
	"github.com/gofrontier-com/ranger/pkg/core"

	"github.com/spf13/cobra"
)

var (
	projectName string
	orgName     string
)

// NewCmdDestroySet creates a command to destroy a set
func NewCmdDestroySet(config *core.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set",
		Short: "Destroy a set",
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := destroy.DestroySet(config, projectName, orgName); err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&projectName, "project-name", "p", "", "Project name")
	cmd.Flags().StringVarP(&orgName, "organisation-name", "o", "", "Organisation name")

	return cmd
}
